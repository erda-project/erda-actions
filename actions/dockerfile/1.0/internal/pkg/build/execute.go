package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/metadata"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/dockerfile"
	"github.com/erda-project/erda-actions/pkg/pack"
)

type Builder interface {
	Build(p *Params, o *OutPut) error
}

type OutPut struct {
	Image string
	Push  bool
}

type Params struct {
	Args         map[string]string
	BuildContext map[string]string
}

// Execute 自定义 dockerfile 构建应用镜像
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

	if cfg.Registry != nil && cfg.Registry.Username != "" {
		if err := docker.Login(cfg.Registry.URL, cfg.Registry.Username, cfg.Registry.Password); err != nil {
			return err
		}
	}

	// docker build & push 业务镜像
	if err := packAndPushImage(cfg); err != nil {
		return err
	}
	return nil
}

func packAndPushImage(cfg conf.Conf) error {
	if cfg.Context != "" {
		if err := os.Chdir(cfg.Context); err != nil {
			return err
		}
	}

	// 判断 dockerfile 是否存在
	if _, err := os.Stat(cfg.Path); err != nil {
		return err
	}

	// 渲染 dockerfile
	if cfg.BuildArgsStr != "" {
		if err := json.Unmarshal([]byte(cfg.BuildArgsStr), &cfg.BuildArgs); err != nil {
			fmt.Printf("failed to unmarshal build_args, :%v\n", err)
			return err
		}

		originalDockerfileContent, err := os.ReadFile(cfg.Path)
		if err != nil {
			return err
		}
		newDockerfileContent := dockerfile.ReplaceOrInsertBuildArgToDockerfile(originalDockerfileContent, cfg.BuildArgs)
		if err = os.WriteFile(cfg.Path, newDockerfileContent, 0644); err != nil {
			return err
		}
	}

	// docker build 业务镜像
	repo := getRepo(cfg)

	buildArgs := map[string]string{
		"NODE_OPTIONS":     fmt.Sprintf("--max_old_space_size=%s", strconv.Itoa(cfg.Memory-100)),
		"PACKAGE_LOCK_DIR": "/.cache_packagejson",
		"DICE_WORKSPACE":   cfg.DiceWorkspace,
	}

	buildContext := cfg.BuildContext
	if buildContext == nil {
		buildContext = make(map[string]string)
	}

	var builder Builder

	isBuildkitEnable, err := strconv.ParseBool(cfg.BuildkitEnable)
	if err != nil {
		fmt.Fprintf(os.Stdout, "unparsed builder keyword: %v\n", err)
		return err
	}

	if isBuildkitEnable {
		builder, err = NewBuildkit(&cfg)
		if err != nil {
			return err
		}
	} else {
		builder, err = NewDocker(&cfg)
		if err != nil {
			return err
		}
	}

	if err := builder.Build(
		&Params{
			Args:         buildArgs,
			BuildContext: buildContext,
		},
		&OutPut{
			Image: repo,
			Push:  true, // push default now.
		},
	); err != nil {
		fmt.Fprintf(os.Stdout, "failed to build: %v\n", err)
		return err
	}

	// upload metadata
	if err := storeMetaFile(&cfg, repo); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully upload metafile\n")
	if cfg.Service != "" { // TODO deprecated
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

	return nil
}

// 生成业务镜像名称
func getRepo(cfg conf.Conf) string {
	// registry url
	registry := cfg.LocalRegistry
	if cfg.Registry != nil && cfg.Registry.URL != "" {
		registry = cfg.Registry.URL
	}
	// image name
	name := cfg.ProjectAppAbbr
	if name == "" {
		name = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	if cfg.Image != nil && cfg.Image.Name != "" {
		name = cfg.Image.Name
	}
	// image tag
	tag := fmt.Sprintf("%s-%v", cfg.TaskName, time.Now().UnixNano())
	if cfg.Image != nil && cfg.Image.Tag != "" {
		tag = cfg.Image.Tag
	}

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(registry), name, tag))
}

func storeMetaFile(cfg *conf.Conf, image string) error {
	meta := apistructs.ActionCallback{
		Metadata: metadata.Metadata{
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
