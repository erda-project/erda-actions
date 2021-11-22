package build

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"

	"github.com/erda-project/erda/pkg/strutil"

	"github.com/erda-project/erda-actions/actions/java/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/pack"
)

// packAndPushAppImage pack and push application image
func packAndPushAppImage(cfg conf.Conf) error {
	// ch workdir
	if err := os.Chdir(cfg.WorkDir); err != nil {
		return err
	}
	// copy assets
	if err := cp(compPrefix, "."); err != nil {
		return err
	}
	if cfg.Assets != "" {
		if err := mustDir(cfg.Assets); err != nil {
			return err
		}
		// copy assets
		if err := cp(cfg.Assets, "./assets"); err != nil {
			return err
		}
	} else {
		// no assets exist, but need create an empty dir
		if err := os.Mkdir("./assets", 0755); err != nil {
			return err
		}
	}
	// check container exist
	ct := fmt.Sprintf("%s/%s", compPrefix, cfg.ContainerType)
	if err := mustDir(ct); err != nil {
		fmt.Fprintf(os.Stdout, "container type: %s not exist", cfg.ContainerType)
		return err
	}

	dockerFilePath := fmt.Sprintf("%s/%s/Dockerfile", compPrefix, cfg.ContainerType)

	dockerCopyCmds := []string{}
	if len(cfg.CopyAssets) > 0 {
		rand := random.New()
		for _, asset := range cfg.CopyAssets {
			idx := strings.Index(asset, ":")
			if idx > 0 {
				source := asset[0:idx]
				if len(asset) == idx {
					// 模板文件路径为空不处理
					fmt.Fprintf(os.Stdout, "invalid asset: %s ", asset)
					continue
				}
				dest := asset[idx+1:]
				absSource := path.Join(cfg.Context, source)
				if strings.Index(source, "/") == 0 {
					//绝对路径
					absSource = source
				}
				sourceTarget := rand.String(6, random.Alphanumeric)
				if err := cp(absSource, sourceTarget); err != nil {
					return err
				}

				dockerCopyCmds = append(dockerCopyCmds, fmt.Sprintf("COPY %s %s", sourceTarget, dest))
			} else {
				if err := cp(path.Join(cfg.Context, asset), "target/"+asset); err != nil {
					return err
				}
			}
		}
	}
	if len(dockerCopyCmds) > 0 {
		dockerFileBytes, _ := ioutil.ReadFile(dockerFilePath)
		for _, cmd := range dockerCopyCmds {
			dockerFileBytes = append(dockerFileBytes, []byte("\n"+cmd)...)
		}
		ioutil.WriteFile(dockerFilePath, dockerFileBytes, os.ModePerm)
	}

	// jar包生成 & docker build 出业务镜像
	repo := getRepo(cfg)

	// 获取用户指定脚本命令
	scriptFile, err := os.Create("pre_start.sh")
	defer scriptFile.Close()
	if err != nil {
		return err
	}

	if cfg.PreStartScript != "" {
		fileInfo, err := os.Stat(cfg.PreStartScript)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return errors.New("is directory, not pre start script file")
		}

		input, err := ioutil.ReadFile(cfg.PreStartScript)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile("pre_start.sh", input, 0644)
		if err != nil {
			return err
		}
	}

	// compose build args
	buildArgs := map[string]string{
		"TARGET":                 "target",
		"MONITOR_AGENT":          cfg.MonitorAgent,
		"SPRING_PROFILES_ACTIVE": cfg.Profile, // TODO: 非 spring 定制,
		"DICE_VERSION":           cfg.DiceVersion,
		"WEB_PATH":               cfg.WebPath,
	}

	if cfg.ContainerVersion != "" {
		buildArgs["CONTAINER_VERSION]"] = cfg.ContainerVersion
	}

	if cfg.PreStartScript != "" {
		buildArgs["SCRIPT_ARGS]"] = cfg.PreStartArgs
	}

	// witch the build method
	if cfg.BuildkitEnable == "true" {
		if err := packWithBuildKit(cfg, repo, buildArgs); err != nil {
			fmt.Fprintf(os.Stdout, "failed to pack with buildkit: %v\n", err)
			return err
		}
	} else {
		if err := packWithDocker(cfg, repo, buildArgs); err != nil {
			fmt.Fprintf(os.Stdout, "failed to pack with docker: %v\n", err)
			return err
		}
	}

	// upload metadata
	if err := storeMetaFile(&cfg, repo); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully upload metafile\n")
	if cfg.Service != "" {
		// TODO Deprecated: 使用 ${java:OUTPUT:image}
		// 写应用镜像信息至 pack-result, 供 release action 读取 & 填充dice.yml
		imageResult := make([]pack.ModuleImage, 0)
		imageResult = append(imageResult, pack.ModuleImage{ModuleName: cfg.Service, Image: repo})
		resultBytes, err := json.MarshalIndent(imageResult, "", "  ")
		if err != nil {
			return err
		}
		if err := filehelper.CreateFile(filepath.Join(cfg.WorkDir, "pack-result"), string(resultBytes), 0644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "successfully write image action: %s\n", repo)
	}

	if err := simpleRun("rm", "-rf", "comp", "target", "assets"); err != nil {
		fmt.Fprintf(os.Stdout, "warning, cleanup failed: %v", err)
	}

	return nil
}

// packWithBuildKit pack and push with buildKit
func packWithBuildKit(cfg conf.Conf, repo string, args map[string]string) error {
	argsSlice := make([]string, 0)

	for k, v := range args {
		argsSlice = append(argsSlice, "--opt", fmt.Sprintf("build-arg:%s=%s", k, v))
	}

	buildCmdArgs := []string{
		"--addr", cfg.BuildkitdAddr,
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
	}

	// append args, e.g. --opt k=v
	buildCmdArgs = append(buildCmdArgs, argsSlice...)

	// append build source and output param.
	buildCmdArgs = append(buildCmdArgs,
		"--local", "context="+cfg.WorkDir,
		"--local", "dockerfile="+fmt.Sprintf("%s/%s", compPrefix, cfg.ContainerType),
		"--output", "type=image,name="+repo+",push=true,registry.insecure=true",
	)

	buildkitCmd := exec.Command("buildctl", buildCmdArgs...)
	fmt.Println(strutil.Join(buildkitCmd.Args, " ", false))

	buildkitCmd.Dir = cfg.WorkDir
	buildkitCmd.Stdout = os.Stdout
	buildkitCmd.Stderr = os.Stderr
	if err := buildkitCmd.Run(); err != nil {
		return err
	}

	return nil
}

// packWithDocker pack and push with docker
func packWithDocker(cfg conf.Conf, repo string, args map[string]string) error {
	argsSlice := make([]string, 0)

	for k, v := range args {
		argsSlice = append(argsSlice, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}

	buildCmdArgs := []string{"build"}
	buildCmdArgs = append(buildCmdArgs, argsSlice...)
	buildCmdArgs = append(buildCmdArgs,
		"--cpu-quota", strconv.FormatFloat(float64(cfg.CPU*100000), 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(cfg.Memory*apistructs.MB), 10),
		"-t", repo,
		"-f", fmt.Sprintf("%s/%s/Dockerfile", compPrefix, cfg.ContainerType), ".")

	packCmd := exec.Command("docker", buildCmdArgs...)

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
	return nil
}

// getRepo compose image full name
func getRepo(cfg conf.Conf) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", cfg.TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}
