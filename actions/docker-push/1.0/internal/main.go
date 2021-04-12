package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/pack"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
)

type Conf struct {
	WorkDir string `env:"WORKDIR"`
	// params
	Image    string `env:"ACTION_IMAGE" required:"true"`
	Username string `env:"ACTION_USERNAME"`
	Password string `env:"ACTION_PASSWORD"`
	From     string `env:"ACTION_FROM"`
	Service  string `env:"ACTION_SERVICE"`
	Pull     bool   `env:"ACTION_PULL"`
	// pipeline 自动注入
	TaskName       string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ProjectAppAbbr string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceOperatorId string `env:"DICE_OPERATOR_ID" default:"terminus"`

	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`
}

func run() error {
	var (
		cfg         Conf
		resultBytes []byte
		err         error
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
		resultBytes, err = pushImage(cfg)
	} else {
		resultBytes, err = pullImage(cfg)
	}

	if err != nil {
		return err
	}

	if err := filehelper.CreateFile(filepath.Join(cfg.WorkDir, "pack-result"), string(resultBytes), 0644); err != nil {
		return err
	}

	return nil
}

func pushImage(cfg Conf) ([]byte, error) {
	fromImage, err := getFrom(cfg)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stdout, "image from: %s\n", fromImage)
	fmt.Fprintf(os.Stdout, "image to: %s\n", cfg.Image)

	if err := simpleRun("docker", "pull", fromImage); err != nil {
		return nil, err
	}

	if err := simpleRun("docker", "tag", fromImage, cfg.Image); err != nil {
		return nil, err
	}

	if cfg.Username != "" {
		// login
		getRegistry := func(image string) string { return strings.Split(image, "/")[0] }
		if err := simpleRun("docker", "login", "-u", cfg.Username, "-p", cfg.Password, getRegistry(cfg.Image)); err != nil {
			return nil, fmt.Errorf("docker login failed, image: %s, username: %s, password: %s, err: %v",
				cfg.Image, cfg.Username, cfg.Password, err)
		}
	}

	if err := simpleRun("docker", "push", cfg.Image); err != nil {
		return nil, err
	}

	imageResult := make([]pack.ModuleImage, 0)
	imageResult = append(imageResult, pack.ModuleImage{ModuleName: cfg.Service, Image: cfg.Image})
	return json.MarshalIndent(imageResult, "", "  ")
}

func pullImage(cfg Conf) ([]byte, error) {
	fmt.Fprintf(os.Stdout, "config: %+v\n", cfg)
	toImage := docker.GetInnerRepoAddr(cfg.ProjectAppAbbr, cfg.DiceOperatorId, cfg.TaskName, cfg.LocalRegistry)
	fmt.Fprintf(os.Stdout, "image from: %s\n", cfg.Image)
	fmt.Fprintf(os.Stdout, "image to: %s\n", toImage)

	if cfg.Username != "" {
		// login
		getRegistry := func(image string) string { return strings.Split(image, "/")[0] }
		if err := simpleRun("docker", "login", "-u", cfg.Username, "-p", cfg.Password, getRegistry(cfg.Image)); err != nil {
			return nil, fmt.Errorf("docker login failed, image: %s, username: %s, password: %s, err: %v",
				cfg.Image, cfg.Username, cfg.Password, err)
		}
	}

	if err := simpleRun("docker", "pull", cfg.Image); err != nil {
		return nil, err
	}

	if err := simpleRun("docker", "tag", cfg.Image, toImage); err != nil {
		return nil, err
	}

	if err := simpleRun("docker", "push", toImage); err != nil {
		return nil, err
	}

	imageResult := make([]pack.ModuleImage, 0)
	imageResult = append(imageResult, pack.ModuleImage{ModuleName: cfg.Service, Image: toImage})
	return json.MarshalIndent(imageResult, "", "  ")
}

func getFrom(cfg Conf) (string, error) {
	v, err := ioutil.ReadFile(cfg.From)
	if err != nil {
		return "", err
	}
	images := make([]pack.ModuleImage, 0)
	if err := json.Unmarshal([]byte(v), &images); err != nil {
		return "", err
	}
	for _, i := range images {
		if cfg.Service == i.ModuleName {
			return i.Image, nil
		}
	}
	return "", fmt.Errorf("not found image of service: %s", cfg.Service)
}

func simpleRun(name string, arg ...string) error {
	fmt.Fprintf(os.Stdout, "Run: %s, %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "docker-push failed, %v", err)
		os.Exit(1)
	}
}
