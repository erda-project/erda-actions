package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
	"github.com/labstack/gommon/random"

	"github.com/erda-project/erda-actions/actions/php/1.0/internal/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/pack"
	"github.com/erda-project/erda-actions/pkg/render"
)

const (
	// components 放置位置
	compPrefix = "/opt/action/comp"
	WebRoot    = "/var/www/html/"
)

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("Chdir %s", cfg.Context))
	err := os.Chdir(cfg.Context)
	if err != nil {
		return err
	}

	if filehelper.CheckExist("composer.json", false) == nil {
		//安装依赖
		fmt.Fprintln(os.Stdout, fmt.Sprintf("install composer dep"))
		if err := runCommand("composer", "config", "-g", "repo.packagist", "composer", "https://mirrors.aliyun.com/composer/"); err != nil {
			return err
		}
		if err := runCommand("composer install"); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("composer.json not found"))
	}

	cfgMap := make(map[string]string)
	cfgMap["APACHE_DOCUMENT_ROOT"] = path.Join(WebRoot, cfg.IndexPath)
	cfgMap["CENTRAL_REGISTRY"] = cfg.CentralRegistry
	cfgMap["PHP_VERSION"] = cfg.PHPVersion
	if err := render.RenderTemplate(compPrefix, cfgMap); err != nil {
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
	if err := os.Chdir(cfg.Context); err != nil {
		return err
	}
	// copy assets
	if err := cp(path.Join(compPrefix, "Dockerfile"), "."); err != nil {
		return err
	}

	// docker build 出业务镜像
	repo := getRepo(cfg)

	if cfg.BuildkitEnable == "true" {
		if err := packWithBuildkit(cfg, repo); err != nil {
			fmt.Fprintf(os.Stdout, "failed to pack with buildkit, %v\n", err)
			return err
		}
	} else {
		if err := packWithDocker(cfg, repo); err != nil {
			fmt.Fprintf(os.Stdout, "failed to pack with docker, %v\n", err)
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

	cleanCmd := exec.Command("rm", "-rf", "comp", "target", "assets")
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "warning, cleanup failed: %v", err)
	}

	return nil
}

func getRepo(cfg conf.Conf) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", cfg.TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}

func packWithDocker(cfg conf.Conf, repo string) error {
	packCmd := exec.Command("docker", "build",
		"--build-arg", fmt.Sprintf("TARGET=%s", "."),
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
	return nil
}

func packWithBuildkit(cfg conf.Conf, repo string) error {
	buildCmdArgs := []string{
		"--addr", cfg.BuildkitdAddr,
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--opt", "build-arg" + fmt.Sprintf("TARGET=%s", "."),
	}

	buildCmdArgs = append(buildCmdArgs,
		"--local", "context="+cfg.Context,
		"--local", "dockerfile="+compPrefix,
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
