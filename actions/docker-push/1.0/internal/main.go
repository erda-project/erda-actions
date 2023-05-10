package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/envconf"
	"github.com/erda-project/erda-actions/pkg/meta"
	"github.com/erda-project/erda-actions/pkg/pack"
)

var getRegistry = func(image string) string { return strings.Split(image, "/")[0] }

type Conf struct {
	WorkDir string `env:"WORKDIR"`
	// params
	Image    string `env:"ACTION_IMAGE" required:"true"`
	Username string `env:"ACTION_USERNAME"`
	Password string `env:"ACTION_PASSWORD"`
	From     string `env:"ACTION_FROM"`
	Service  string `env:"ACTION_SERVICE"`
	Pull     bool   `env:"ACTION_PULL"`
	Insecure bool   `env:"ACTION_INSECURE" default:"true"`

	// pipeline 自动注入
	TaskName       string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ProjectAppAbbr string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceOperatorId string `env:"DICE_OPERATOR_ID" default:"terminus"`

	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`

	// BuildKit params
	BuildkitEnable string `env:"BUILDKIT_ENABLE"`
	BuildkitdAddr  string `env:"BUILDKITD_ADDR" default:"tcp://buildkitd.default.svc.cluster.local:1234"`
}

func run() error {
	var (
		cfg         Conf
		resultBytes []byte
		err         error
		image       string
	)
	if err := envconf.Load(&cfg); err != nil {
		return err
	}

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	if !cfg.Pull {
		image = cfg.Image
		resultBytes, err = pushImage(cfg)
	} else {
		image = docker.GetInnerRepoAddr(cfg.ProjectAppAbbr, cfg.DiceOperatorId, cfg.TaskName, cfg.LocalRegistry)
		resultBytes, err = pullImage(cfg, image)
	}

	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to process image: %v\n", err)
		return err
	}

	collector := meta.NewResultMetaCollector()
	collector.Add("image", image)
	logrus.Infof("successfully upload metafile")

	if err := filehelper.CreateFile(filepath.Join(cfg.WorkDir, "pack-result"), string(resultBytes), 0644); err != nil {
		return err
	}
	logrus.Infof("successfully write pack-result")
	return nil
}

func pushImage(cfg Conf) ([]byte, error) {
	fromImage, err := getFrom(cfg)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stdout, "image from: %s\n", fromImage)
	fmt.Fprintf(os.Stdout, "image to: %s\n", cfg.Image)

	if cfg.Username != "" {
		// login
		if err := docker.Login(getRegistry(cfg.Image), cfg.Username, cfg.Password); err != nil {
			return nil, fmt.Errorf("failed to login, error: %v", err)
		}
	}

	if err := reTag(fromImage, cfg.Image, cfg.Insecure); err != nil {
		return nil, err
	}

	imageResult := make([]pack.ModuleImage, 0)
	imageResult = append(imageResult, pack.ModuleImage{ModuleName: cfg.Service, Image: cfg.Image})
	return json.MarshalIndent(imageResult, "", "  ")
}

func pullImage(cfg Conf, toImage string) ([]byte, error) {
	fmt.Fprintf(os.Stdout, "image from: %s\n", cfg.Image)
	fmt.Fprintf(os.Stdout, "image to: %s\n", toImage)

	if cfg.Username != "" {
		// login
		if err := docker.Login(getRegistry(cfg.Image), cfg.Username, cfg.Password); err != nil {
			return nil, fmt.Errorf("failed to login, error: %v", err)
		}
	}

	if err := reTag(cfg.Image, toImage, cfg.Insecure); err != nil {
		return nil, err
	}

	imageResult := make([]pack.ModuleImage, 0)
	imageResult = append(imageResult, pack.ModuleImage{ModuleName: cfg.Service, Image: toImage})
	return json.MarshalIndent(imageResult, "", "  ")
}

func getFrom(cfg Conf) (string, error) {
	v, err := os.ReadFile(cfg.From)
	if err != nil {
		return "", err
	}
	images := make([]pack.ModuleImage, 0)
	if err := json.Unmarshal(v, &images); err != nil {
		return "", err
	}
	for _, i := range images {
		if cfg.Service == i.ModuleName {
			return i.Image, nil
		}
	}
	return "", fmt.Errorf("not found image of service: %s", cfg.Service)
}

func reTag(from, target string, insecure bool) error {
	fn := func(arg ...string) error {
		fmt.Fprintf(os.Stdout, "Run: gcrane, %v\n", arg)
		cmd := exec.Command("gcrane", arg...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	imageFile := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	args := []string{
		fmt.Sprintf("--insecure=%s", strconv.FormatBool(insecure)),
	}

	pullArgs := append([]string{"pull", from, imageFile}, args...)
	if err := fn(pullArgs...); err != nil {
		return err
	}

	pushArgs := append([]string{"push", imageFile, target}, args...)
	if err := fn(pushArgs...); err != nil {
		return err
	}

	return nil
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "docker-push failed, %v", err)
		os.Exit(1)
	}
}
