package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/random"

	"github.com/erda-project/erda-actions/actions/golang/1.0/internal/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/pack"
	"github.com/erda-project/erda-actions/pkg/render"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	// components 放置位置
	compPrefix = "/opt/action/comp"
)

var goPATH = os.Getenv("GOPATH")
var cfg conf.Conf

func Execute() error {

	envconf.MustLoad(&cfg)

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	buildInfo, err := getBuildInfo(cfg)
	if err != nil {
		return err
	}
	if cfg.Target == "" {
		cfg.Target = filepath.Base(buildInfo.Package)
	}

	if cfg.Service == "" {
		cfg.Service = filepath.Base(buildInfo.Package)
	}

	buildPath := buildInfo.BuildPath
	fmt.Fprintln(os.Stdout, fmt.Sprintf("build path :%s", buildPath))
	if buildInfo.BuildType != GoMODBuild {
		// 不是gomod 需要把代码放到gopath下
		if err := runCommand("mkdir", "-p", buildPath); err != nil {
			return err
		}
		if err := runCommand("mv", cfg.Context+"/*", buildPath); err != nil {
			return err
		}
	}
	if err := os.Chdir(buildPath); err != nil {
		return err
	}

	err = os.Mkdir(path.Join(cfg.WorkDir, "assets"), os.ModePerm)
	if err != nil {
		return err
	}

	cfgMap := make(map[string]string)
	cfgMap["TARGET"] = path.Join(cfg.WorkDir, "target")
	if cfg.Assets != nil && len(cfg.Assets) > 0 {
		for _, assetPath := range cfg.Assets {
			pathList := strings.Split(assetPath, ":")
			if len(pathList) == 2 {
				destPath := path.Join(cfg.WorkDir, "assets", pathList[1])
				err = os.MkdirAll(path.Dir(destPath), os.ModePerm)
				if err != nil {
					return err
				}
				fmt.Fprintf(os.Stdout, "copy asset %s => %s\n", pathList[0], destPath)
				err := runCommand("cp", "-r", pathList[0], destPath)
				if err != nil {
					return err
				}
			} else {
				fmt.Fprintf(os.Stdout, "copy asset %s => %s\n", assetPath, path.Join(cfg.WorkDir, "assets"))
				err := runCommand("cp", "--parents", "-r", assetPath, path.Join(cfg.WorkDir, "assets"))
				if err != nil {
					return err
				}
			}
		}
	}
	if err := render.RenderTemplate(compPrefix, cfgMap); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("begin build target"))
	if err := runCommand(cfg.Command); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, fmt.Sprintf("build target success"))

	if _, err := os.Stat(cfg.Target); err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("target file %s not exist", cfg.Target))
		}
		return err
	}
	if err := cp(cfg.Target, path.Join(cfg.WorkDir, "target")); err != nil {
		return err
	}

	err = packAndPushAppImage(cfg)
	if err != nil {
		return err
	}
	return nil
}

func runCommand(cmd ...string) error {
	c := strings.Join(cmd, " ")
	buildCmd := exec.Command("/bin/bash", "-c", c)

	if cfg.GoProxy != "" {
		buildCmd.Env = os.Environ()
		buildCmd.Env = append(buildCmd.Env, "GOPROXY="+cfg.GoProxy)
	}

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}
	return nil
}

// docker build & docker push 业务镜像
func packAndPushAppImage(cfg conf.Conf) error {
	// 切换工作目录
	if err := os.Chdir(cfg.WorkDir); err != nil {
		return err
	}
	// copy assets
	if err := cp(path.Join(compPrefix, "Dockerfile"), "."); err != nil {
		return err
	}

	// docker build 出业务镜像
	repo := getRepo(cfg)

	if cfg.BuildkitEnable == "true" {
		if err := packWithBuildkit(repo); err != nil {
			return err
		}
	} else {
		if err := packWithDocker(repo); err != nil {
			return err
		}
	}
	// upload metadata
	if err := storeMetaFile(&cfg, repo); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully upload metafile\n")

	if err := storePackResult(cfg, repo); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully write pack-result\n")

	return nil
}

func packWithDocker(repo string) error {
	packCmd := exec.Command("docker", "build",
		"--build-arg", fmt.Sprintf("TARGET=%s", "target"),
		"--cpu-quota", strconv.FormatFloat(cfg.CPU*100000, 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(cfg.Memory*apistructs.MB), 10),
		"-t", repo,
		"-f", fmt.Sprintf("Dockerfile"), ".")

	fmt.Fprintf(os.Stdout, "packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully build app image: %s\n", repo)

	// docker push 业务镜像至集群 registry
	if err := docker.PushByCmd(repo, ""); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully push app image: %s\n", repo)
	return  nil
}

func packWithBuildkit(repo string) error {
	packCmd := exec.Command("buildctl",
		"--addr",
		"tcp://buildkitd.default.svc.cluster.local:1234",
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--opt", "build-arg:TARGET=target",
		"--local", "context=" + cfg.WorkDir,
		"--local", "dockerfile=" + cfg.WorkDir,
		"--output", "type=image,name=" + repo + ",push=true,registry.insecure=true")

	fmt.Fprintf(os.Stdout, "packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully build app image: %s\n", repo)
	return  nil
}

func getRepo(cfg conf.Conf) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", cfg.TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}

func storeMetaFile(cfg *conf.Conf, image string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "image",
				Value: image,
			},
		},
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(cfg.MetaFile, string(b), 0644); err != nil {
		return errors.New("write file:metafile failed")
	}
	return nil
}

func storePackResult(cfg conf.Conf, repo string) error {
	imageResult := make([]pack.ModuleImage, 0)
	moduleName := cfg.Service
	if moduleName == "" {
		moduleName = cfg.TaskName
	}
	imageResult = append(imageResult, pack.ModuleImage{ModuleName: moduleName, Image: repo})
	resultBytes, err := json.MarshalIndent(imageResult, "", "  ")
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(filepath.Join(cfg.WorkDir, "pack-result"), string(resultBytes), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully write image action: %s\n", repo)
	return nil
}

func cp(srcDir, destDir string) error {
	fmt.Fprintf(os.Stdout, "copying %s to %s\n", srcDir, destDir)
	cpCmd := exec.Command("cp", "-r", srcDir, destDir)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	return cpCmd.Run()
}

func getBuildInfo(cfg conf.Conf) (*GoBuild, error) {
	goModPath := path.Join(cfg.Context, "go.mod")
	// go mod 返回代码原路径
	if filehelper.CheckExist(goModPath, false) == nil {
		fmt.Fprintf(os.Stdout, "go mod detect\n")
		goModBytes, err := ioutil.ReadFile(goModPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to read go mod file %s err:%s", goModPath, err.Error()))
		}
		return &GoBuild{
			BuildPath: cfg.Context,
			BuildType: GoMODBuild,
			Package:   ModulePath(goModBytes),
		}, nil
	}
	goVendorPath := path.Join(cfg.Context, "vendor", "vendor.json")
	if filehelper.CheckExist(goVendorPath, false) == nil {
		vendorBytes, err := ioutil.ReadFile(goVendorPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to read go verdor file %s err:%s", goVendorPath, err.Error()))
		}
		var vendorConfig GoVendorConfig
		err = json.Unmarshal(vendorBytes, &vendorConfig)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to parse go verdor file %s err:%s", goVendorPath, err.Error()))
		}
		fmt.Fprintf(os.Stdout, fmt.Sprintf("go vendor detected package:%s\n", vendorConfig.RootPath))
		return &GoBuild{
			BuildPath: goPATH + vendorConfig.RootPath,
			BuildType: GoVendorBuild,
			Package:   vendorConfig.RootPath,
		}, nil
	}

	// 无法探测 使用指定的package
	if cfg.Package != "" {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("use config package:%s\n", cfg.Package))
		return &GoBuild{
			BuildPath: goPATH + cfg.Package,
			BuildType: OtherBuild,
			Package:   cfg.Package,
		}, nil
	}
	return nil, errors.New("failed to get go package name")
}
