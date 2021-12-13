package build

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/render"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/gitbook/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	compPrefix = "/opt/action/comp"
	nginxConf  = "nginx.conf.template"
)

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)
	fmt.Fprintln(os.Stdout, "sucessfully loaded action config")

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	cfgMap := make(map[string]string)
	cfgMap["CENTRAL_REGISTRY"] = cfg.CentralRegistry
	if err := render.RenderTemplate(compPrefix, cfgMap); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "successfully replaced action placeholder")

	// 编译打包应用
	if err := build(cfg); err != nil {
		return err
	}

	// docker build & push 业务镜像
	if err := packAndPushImage(cfg); err != nil {
		return err
	}
	return nil
}

func build(cfg conf.Conf) error {
	// 切换工作目录
	if cfg.Context != "" {
		fmt.Fprintf(os.Stdout, "change workding directory to: %s\n", cfg.Context)
		if err := os.Chdir(cfg.Context); err != nil {
			return err
		}
	}

	buildCmd := exec.Command("gitbook", "build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	fmt.Fprintf(os.Stdout, "buildCmd: %v\n", buildCmd.Args)
	if err := buildCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "built successfully")
	destPath := "_book"

	// 校验构建完成的目标目录是否存在, eg: public
	if _, err := os.Stat(destPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(cfg.WorkDir, destPath), 0755); err != nil {
		return err
	}
	if err := cp(destPath, cfg.WorkDir); err != nil {
		return err
	}

	confFile := fmt.Sprintf("%s/%s", compPrefix, nginxConf)

	if err := cp(confFile, cfg.WorkDir); err != nil {
		return err
	}
	if err := cp(fmt.Sprintf("%s/%s", compPrefix, "Dockerfile"), cfg.WorkDir); err != nil {
		return err
	}

	return nil
}

func packAndPushImage(cfg conf.Conf) error {
	if err := os.Chdir(cfg.WorkDir); err != nil {
		return err
	}
	if err := cp(compPrefix, "."); err != nil {
		return err
	}

	// docker build 业务镜像
	repo := getRepo(cfg)
	if cfg.BuildkitEnable == "true" {
		if err := packWithBuildkit(repo, cfg); err != nil {
			return err
		}
	} else {
		if err := packWithDocker(repo, cfg); err != nil {
			return err
		}
	}
	// upload metadata
	if err := storeMetaFile(&cfg, repo); err != nil {
		return err
	}

	cleanCmd := exec.Command("rm", "-rf", compPrefix)
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "warning, cleanup failed: %v", err)
	}

	return nil
}

func packWithDocker(repo string, cfg conf.Conf) error {
	packCmd := exec.Command("docker", "build",
		"--build-arg", fmt.Sprintf("DICE_VERSION=%s", cfg.DiceVersion),
		"--cpu-quota", strconv.FormatFloat(float64(cfg.CPU*100000), 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(cfg.Memory*apistructs.MB), 10),
		"-t", repo,
		".")
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

func packWithBuildkit(repo string, cfg conf.Conf) error {
	packCmd := exec.Command("buildctl",
		"--addr",
		"tcp://buildkitd.default.svc.cluster.local:1234",
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--opt", "build-arg:" + fmt.Sprintf("DICE_VERSION=%s", cfg.DiceVersion),
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


// 生成业务镜像名称
func getRepo(cfg conf.Conf) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", cfg.TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}

func cp(src, dst string, fileType ...string) error {
	var cpCmd *exec.Cmd
	if len(fileType) > 0 {
		ft := fileType[0]
		cpCmd = exec.Command("find", src, fmt.Sprintf("*.%v", ft), "-exec", "cp", "{}", dst, "\\;")
	} else {
		cpCmd = exec.Command("cp", "-r", src, dst)
	}
	fmt.Fprintf(os.Stdout, "cpCmd: %v\n", cpCmd.Args)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	return cpCmd.Run()
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
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}
