package run

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/dlog"
	"github.com/erda-project/erda/pkg/strutil"
)

func executeCmd() (string, error) {
	cmds := strutil.Split(conf.UserConf().Cmd, " ")
	deployCmd := exec.Command(cmds[0], cmds[1:]...)
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	deployCmd.Dir = conf.UserConf().Workdir
	dlog.Printf("will execute specified cmd: %s\n", deployCmd.String())

	if err := deployCmd.Run(); err != nil {
		return "", errors.Errorf("failed to execute specified cmd, err: %v", err)
	}

	return "", nil
}
