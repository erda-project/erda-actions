package build

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/erda-project/erda-actions/actions/git-push/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

// Execute 推送 npm library 至远程 registry
func Execute() error {
	// 1. git init .
	// 2. git config --global user.email "dice@terminus.io"
	// 3. git config --global user.name "dice"
	// 4. git add .
	// 5. git commit -m "init templates"
	// 6. git add origin remoteUrl
	// 7. git push origin master
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

	gitPushCmd := exec.Command("/bin/bash", "/opt/action/git-push.sh", cfg.RemoteUrl)
	gitPushCmd.Stdout = os.Stdout
	gitPushCmd.Stderr = os.Stderr
	if err := gitPushCmd.Run(); err != nil {
		return err
	}

	return nil
}
