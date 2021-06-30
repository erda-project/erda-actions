package build

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/js-deploy/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
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

// Execute 推送 npm library 至远程 registry
func Execute() error {
	// 1. npm config set registry xxx
	// 2. npm login
	// 3. npm publish
	var cfg conf.Conf
	envconf.MustLoad(&cfg)
	fmt.Fprintln(os.Stdout, "sucessfully loaded action config")

	// 切换工作目录
	if cfg.Context != "" {
		fmt.Fprintf(os.Stdout, "change workding directory to: %s\n", cfg.Context)
		if err := os.Chdir(cfg.Context); err != nil {
			return err
		}
	}

	// Npm Registry 地址须以 / 结尾
	if !strings.HasSuffix(cfg.NpmRegistry, "/") {
		cfg.NpmRegistry = cfg.NpmRegistry + "/"
	}

	// npm config
	fmt.Fprintf(os.Stdout, "setting npm registry: %s\n", cfg.NpmRegistry)
	npmConfigCmd := exec.Command("npm", "config", "set", "@terminus:registry", cfg.NpmRegistry)
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

	// npm publish
	fmt.Fprintf(os.Stdout, "starting publish npm package to registry %s\n", cfg.NpmRegistry)
	npmPublishCmd := exec.Command("npm", "publish", "--registry", cfg.NpmRegistry)
	npmPublishCmd.Stdout = os.Stdout
	npmPublishCmd.Stderr = os.Stderr
	if err := npmPublishCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "finished publish npm package to registry")

	return nil
}
