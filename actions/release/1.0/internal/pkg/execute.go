package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
)

var errReleaseTypeCheck = errors.New(`一个release action只能发布一种操作系统类型的移动应用，如果有多种操作系统类型，请拆分成多个release acction`)

const (
	andriodExt   = ".apk"
	iosExt       = ".ipa"
	distFilePath = "/tmp/dist.zip"
)

// Execute release action 执行逻辑
func Execute() error {
	logrus.SetOutput(os.Stdout)

	var cfg conf.Conf
	if err := initEnv(&cfg); err != nil {
		return err
	}

	// check application mode
	app, err := GetApp(cfg.AppID, cfg)
	if err != nil {
		return err
	}
	// project level application not support release action
	if app.Mode == string(apistructs.ApplicationModeProjectService) {
		return fmt.Errorf("project level application not support release action")
	}

	// generate release create request
	req := genReleaseRequest(&cfg)
	storage, err := parseURL(cfg.PipelineStorageURL)
	if err != nil {
		return err
	}
	//image和services是release的两种模式，一种是用户传入image，然后release做些其他处理，
	//一种就是传入service，然后release根据配置进行构建出imageAddr, 然后把imageAddr设置进images属性中
	//这里就是新加的service，根据service构建，然后塞入image中
	if isServiceMode(&cfg) && !DockerBuildPushAndSetImages(&cfg) {
		return errors.Errorf("release build and push error")
	}

	// 上传 release files & 填充 release resources
	if cfg.ReleaseFiles != "" {
		releaseResources, err := UploadFiles(cfg.ReleaseFiles, storage)
		if err != nil {
			return err
		}
		req.Resources = releaseResources
		fmt.Println(fmt.Sprintf("uploaded release resources: %+v", releaseResources))
	}

	if cfg.ReleaseMobile != nil {
		var typeCounts int
		// 移动应用文件资源
		// 一次release只能release一种类型的移动应用
		for _, appFilePath := range cfg.ReleaseMobile.Files {
			ext := filepath.Ext(appFilePath)
			if ext == andriodExt || ext == iosExt || (ext == "" && filepath.Base(appFilePath) == "dist") {
				typeCounts++
			}
		}
		if typeCounts > 1 {
			return errReleaseTypeCheck
		}

		for _, appFilePath := range cfg.ReleaseMobile.Files {
			var resourceType apistructs.ResourceType
			ext := filepath.Ext(appFilePath)
			if ext == ".apk" {
				resourceType = apistructs.ResourceTypeAndroid
			} else if ext == ".aab" {
				resourceType = apistructs.ResourceTypeAndroidAppBundle
			} else if ext == ".ipa" {
				resourceType = apistructs.ResourceTypeIOS
			} else if ext == "" && filepath.Base(appFilePath) == "dist" {
				resourceType = apistructs.ResourceTypeH5
			} else {
				resourceType = apistructs.ResourceTypeDataSet
			}

			var mobileFileUploadResult *apistructs.FileDownloadFailResponse
			if resourceType != apistructs.ResourceTypeH5 {
				mobileFileUploadResult, err = UploadFileNew(appFilePath, cfg)
				if err != nil {
					return err
				}
			}
			logoTmpPath := "/tmp/logo.jpg"
			meta := map[string]interface{}{}
			meta["logo"] = ""
			if resourceType == apistructs.ResourceTypeAndroid {
				info, err := GetAndroidAppInfo(appFilePath)
				if err != nil {
					return err
				}
				versionCodeStr := strconv.FormatInt(int64(info.VersionCode), 10)
				version := info.Version
				meta["packageName"] = info.PackageName
				meta["version"] = version
				meta["buildID"] = versionCodeStr
				meta["displayName"] = info.Version
				if info.Icon != nil {
					err = SaveImageToFile(info.Icon, logoTmpPath)
					if err != nil {
						logrus.Errorf("error encode jpeg icon %s %v", appFilePath, err)
					} else {
						logoUploadResult, err := UploadFileNew(logoTmpPath, cfg)
						if err != nil {
							return err
						}
						meta["logo"] = logoUploadResult.Data.DownloadURL
					}
				}
				if cfg.ReleaseMobile.Version == "" {
					cfg.ReleaseMobile.Version = version
				}
				req.Version = version
			}
			// TODO Change get information from configuration to extract from abb file
			if resourceType == apistructs.ResourceTypeAndroidAppBundle {
				// info, err := GetAndroidAppBundleInfo(appFilePath)
				// if err != nil {
				// 	return err
				// }
				if cfg.AABInfo.PackageName == "" {
					return errors.Errorf("aab's package name is empty")
				}
				versionCode := cfg.PipelineID
				if cfg.AABInfo.VersionCode != "" {
					versionCode = strutil.String(cfg.AABInfo.VersionCode)
				}
				meta["packageName"] = cfg.AABInfo.PackageName
				meta["version"] = cfg.AABInfo.VersionName
				meta["buildID"] = versionCode
				meta["displayName"] = cfg.AABInfo.VersionName
				if cfg.ReleaseMobile.Version == "" {
					cfg.ReleaseMobile.Version = strutil.String(cfg.AABInfo.VersionName)
				}
				req.Version = strutil.String(cfg.AABInfo.VersionName)
			}

			if resourceType == apistructs.ResourceTypeIOS {
				info, err := GetIOSAppInfo(appFilePath)
				if err != nil {
					return err
				}
				version := info.Version
				meta["bundleId"] = info.BundleId
				meta["build"] = info.Build
				meta["buildID"] = info.Build
				meta["version"] = version
				meta["packageName"] = info.Name
				meta["displayName"] = info.Name
				meta["appStoreURL"] = getAppStoreURL(info.BundleId)
				if info.Icon != nil {
					err = SaveImageToFile(info.Icon, logoTmpPath)
					if err != nil {
						logrus.Errorf("error encode jpeg icon %s %v", appFilePath, err)
					} else {
						logoUploadResult, err := UploadFileNew(logoTmpPath, cfg)
						if err != nil {
							return err
						}
						meta["logo"] = logoUploadResult.Data.DownloadURL
					}
				}
				installPlistContent := GenerateInstallPlist(info, mobileFileUploadResult.Data.DownloadURL)
				plistFile := "/tmp/install.plist"
				err = ioutil.WriteFile(plistFile, []byte(installPlistContent), os.ModePerm)
				if err != nil {
					return err
				}
				plistFileUpploadResult, err := UploadFileNew(plistFile, cfg)
				if err != nil {
					return err
				}
				meta["installPlist"] = plistFileUpploadResult.Data.DownloadURL
				if cfg.ReleaseMobile.Version == "" {
					cfg.ReleaseMobile.Version = version
				}
				req.Version = version
			}
			if resourceType == apistructs.ResourceTypeH5 {
				var h5VersionInfo apistructs.H5VersionInfo
				f, err := ioutil.ReadFile(appFilePath + "/mobileBuild.cfg")
				if err != nil {
					return errors.Errorf("Get h5 version info err: %v", err)
				}
				if err := json.Unmarshal(f, &h5VersionInfo); err != nil {
					return err
				}
				if h5VersionInfo.BuildID == "" {
					h5VersionInfo.BuildID = cfg.PipelineID
				}
				if err := Zip(appFilePath, distFilePath); err != nil {
					return err
				}
				mobileFileUploadResult, err = UploadFileNew(distFilePath, cfg)
				if err != nil {
					return err
				}
				fmt.Println(fmt.Sprintf("H5VersionInfo is %v", h5VersionInfo))

				version := h5VersionInfo.Version
				buildID := h5VersionInfo.BuildID
				if buildID == "" {
					buildID = cfg.PipelineID
				}
				meta["version"] = version
				meta["buildID"] = buildID
				meta["packageName"] = h5VersionInfo.PackageName
				vinfo, err := json.Marshal(&h5VersionInfo)
				if err != nil {
					return err
				}
				meta["h5VersionInfo"] = vinfo
				req.Version = version
			}

			meta["byteSize"] = mobileFileUploadResult.Data.ByteSize
			meta["fileId"] = mobileFileUploadResult.Data.ID
			req.Version = string(resourceType) + "-" + req.Version + "-" + time.Now().Format("20060102150405")

			req.Resources = append(req.Resources, apistructs.ReleaseResource{
				Type: resourceType,
				Name: mobileFileUploadResult.Data.DisplayName,
				URL:  mobileFileUploadResult.Data.DownloadURL,
				Meta: meta,
			})
		}
	}
	// 填充 dice.yml(合并对应环境dice.yml & 填充dice.yml镜像)
	diceYml, err := fillDiceYml(&cfg, storage)
	if err != nil && cfg.ReleaseMobile == nil {
		return err
	}
	req.Dice = diceYml
	fmt.Println(fmt.Sprintf("composed & filled dice.yml: %v", req.Dice))

	// migration sql release
	migrationReleaseID, err := migration(&cfg)
	if err != nil {
		return err
	}
	// 判断resource是否为空
	if len(req.Resources) == 0 {
		req.Resources = []apistructs.ReleaseResource{}
	}
	if migrationReleaseID != "" {
		migResource := apistructs.ReleaseResource{
			Type: apistructs.ResourceTypeMigration,
			Name: apistructs.MigrationResourceKey,
			URL:  migrationReleaseID,
		}
		req.Resources = append(req.Resources, migResource)
	}

	// push release to dicehub
	releaseID, err := pushRelease(cfg, req)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("releaseId: %s", releaseID))

	// create dicehub_release file in nfs, store releaseID
	if err = ioutil.WriteFile("dicehub_release", []byte(releaseID), 0644); err != nil {
		return errors.Wrap(err, "failed to store release id")
	}
	// write metafile
	metaInfos := make([]apistructs.MetadataField, 0, 1)
	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  apistructs.ActionCallbackReleaseID,
		Value: releaseID,
		Type:  apistructs.ActionCallbackTypeLink,
	})
	metaByte, _ := json.Marshal(apistructs.ActionCallback{Metadata: metaInfos})
	if err = filehelper.CreateFile(cfg.Metafile, string(metaByte), 0644); err != nil {
		logrus.Warnf("failed to write metafile, %v", err)
	}

	return nil
}

func isServiceMode(cfg *conf.Conf) bool {
	return cfg.Services != nil && cfg.Images == nil
}

func simpleRun(name string, arg ...string) error {
	fmt.Fprintf(os.Stdout, "Run: %s, %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommand(cmd string) error {
	command := exec.Command("/bin/bash", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

func initEnv(cfg *conf.Conf) error {
	envconf.MustLoad(cfg)

	if cfg.ServicesStr != "" {
		services := make(map[string]conf.Service)
		if err := json.Unmarshal([]byte(cfg.ServicesStr), &services); err != nil {
			return err
		}
		cfg.Services = services
	}

	if cfg.LabelStr != "" {
		if err := json.Unmarshal([]byte(cfg.LabelStr), &cfg.Labels); err != nil {
			return err
		}
	}
	if cfg.ReplacementImageStr != "" {
		if err := json.Unmarshal([]byte(cfg.ReplacementImageStr), &cfg.ReplacementImages); err != nil {
			return err
		}
	}
	if cfg.ImageStr != "" {
		if err := json.Unmarshal([]byte(cfg.ImageStr), &cfg.Images); err != nil {
			return err
		}
	}

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	if cfg.AABInfoStr != "" {
		if err := json.Unmarshal([]byte(cfg.AABInfoStr), &cfg.AABInfo); err != nil {
			return err
		}
	}

	return nil
}

func GetApp(id int64, conf conf.Conf) (*apistructs.ApplicationDTO, error) {

	var resp apistructs.ApplicationFetchResponse

	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/applications/%d", id)).
		Header("Authorization", conf.CiOpenapiToken).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, status-code: %d, content-type: %s", response.StatusCode(), response.ResponseHeader("Content-Type")))
	}

	if !resp.Success {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, error code: %s, error message: %s", resp.Error.Code, resp.Error.Msg))
	}

	return &resp.Data, nil
}
