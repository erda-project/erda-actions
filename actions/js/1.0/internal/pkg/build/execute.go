package build

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/js/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/pack"
	"github.com/erda-project/erda-actions/pkg/render"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

const (
	compPrefix = "/opt/action/comp"

	containerSPA  = "spa"
	containerHerd = "herd"

	nginxConf = "nginx.conf.template"
)

// NpmLoginReq npm login 请求
type NpmLoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// NpmLoginResp npm login 响应
type NpmLoginResp struct {
	Rev   string `json:"rev"`
	ID    string `json:"id"`
	OK    string `json:"ok"`
	Token string `json:"token"`
}

// Execute 构建 js 应用，生成 js 应用镜像
func Execute() error {
	// 1. 参数解析, 选择构建 dockerfile
	// 2. 填充 dockerfile
	// 3. 编译 js 应用
	// 4. 构建 js 镜像, 推送至 registry
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
	cfgMap["DESTDIR"] = cfg.DestDir
	if cfg.NpmRegistry != "" && strings.HasPrefix(cfg.NpmRegistry, "http") {
		cfgMap["NPM_REGISTRY"] = cfg.NpmRegistry
	}
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

	// 兼容离线打包
	if cfg.NpmUsername != "" && cfg.NpmPassword != "" {
		// npm config
		fmt.Fprintf(os.Stdout, "setting npm registry: %s\n", cfg.NpmRegistry)
		npmConfigCmd := exec.Command("npm", "config", "set", "registry", cfg.NpmRegistry)
		npmConfigCmd.Stdout = os.Stdout
		npmConfigCmd.Stderr = os.Stderr
		if err := npmConfigCmd.Run(); err != nil {
			return err
		}

		npmConfigCmd = exec.Command("npm", "config", "set", "@terminus:registry", cfg.NpmRegistry)
		npmConfigCmd.Stdout = os.Stdout
		npmConfigCmd.Stderr = os.Stderr
		if err := npmConfigCmd.Run(); err != nil {
			return err
		}

		u, err := url.Parse(cfg.NpmRegistry)
		if err != nil {
			return err
		}
		npmLoginReq := NpmLoginReq{
			Name:     cfg.NpmUsername,
			Password: cfg.NpmPassword,
		}
		fmt.Fprintln(os.Stdout, "npm login...")
		var npmLoginResp NpmLoginResp
		var client *httpclient.HTTPClient
		if strings.HasPrefix(cfg.NpmRegistry, "https") {
			client = httpclient.New(httpclient.WithHTTPS())
		} else {
			client = httpclient.New()
		}
		resp, err := client.Put(u.Host).
			Path(fmt.Sprintf("%s-/user/org.couchdb.user:%s", u.Path, cfg.NpmUsername)).
			JSONBody(npmLoginReq).Do().JSON(&npmLoginResp)
		if err != nil {
			return err
		}
		if !resp.IsOK() {
			fmt.Fprintf(os.Stdout, "npm login failed")
			return errors.Errorf("npm login failed")
		}
		fmt.Fprintln(os.Stdout, "npm login success")

		fmt.Fprintln(os.Stdout, "setting token to ~/.npmrc...")
		npmConfigCmd = exec.Command("npm", "config", "set", fmt.Sprintf("//%s%s:_authToken=%s", u.Host, u.Path, npmLoginResp.Token))
		npmConfigCmd.Stdout = os.Stdout
		npmConfigCmd.Stderr = os.Stderr
		if err := npmConfigCmd.Run(); err != nil {
			return err
		}
	}

	// 下载应用依赖
	dependencyCmd := exec.Command("/bin/bash", "-c", cfg.DependencyCmd)
	dependencyCmd.Stdout = os.Stdout
	dependencyCmd.Stderr = os.Stderr
	if err := dependencyCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "successfully downloaded dependencies")

	buildCmd := exec.Command("/bin/bash", "-c", cfg.BuildCmd)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	fmt.Fprintf(os.Stdout, "buildCmd: %v\n", buildCmd.Args)
	if err := buildCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "built successfully")

	// 校验构建完成的目标目录是否存在, eg: public
	if _, err := os.Stat(cfg.DestDir); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(cfg.WorkDir, cfg.DestDir), 0755); err != nil {
		return err
	}
	if err := cp(cfg.DestDir, filepath.Join(cfg.WorkDir, filepath.Dir(cfg.DestDir))); err != nil {
		return err
	}

	switch cfg.ContainerType {
	// spa 应用时，检查 nginx.conf.template 是否存在
	case containerSPA:
		// spa 类型容器需要 nginx.conf.template
		if _, err := os.Stat(nginxConf); err != nil {
			return err
		}
		if err := cp(nginxConf, cfg.WorkDir); err != nil {
			return err
		}
	case containerHerd:
		// herd 情况下拷贝应用所有文件
		if err := cp(".", cfg.WorkDir); err != nil {
			return err
		}
	default:
		return errors.New("invalid param container_type")
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
	packCmd := exec.Command("docker", "build",
		"--build-arg", fmt.Sprintf("DSTDIR=%s", cfg.DestDir),
		"--build-arg", fmt.Sprintf("DICE_VERSION=%s", cfg.DiceVersion),
		"--cpu-quota", strconv.FormatFloat(float64(cfg.CPU*100000), 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(cfg.Memory*apistructs.MB), 10),
		"-t", repo,
		"-f", fmt.Sprintf("%s/%s/Dockerfile", filepath.Base(compPrefix), cfg.ContainerType),
		".")
	fmt.Fprintf(os.Stdout, "packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully build app image: %s\n", repo)

	// docker push 业务镜像至集群 registry
	appPushCmd := exec.Command("docker", "push", repo)
	appPushCmd.Stdout = os.Stdout
	appPushCmd.Stderr = os.Stderr
	if err := appPushCmd.Run(); err != nil {
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

	cleanCmd := exec.Command("rm", "-rf", compPrefix)
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "warning, cleanup failed: %v", err)
	}

	return nil
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
