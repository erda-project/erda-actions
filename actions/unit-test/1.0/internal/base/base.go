package base

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
)

const (
	ExtraOsName    = "os.name"
	ExtraOsArch    = "os.arch"
	ExtraOsVersion = "os.ersion"
	Java           = "java"
	Js             = "js"
	Golang         = "golang"

	// ut-action exexute status, finished always.
	UTSatus = "FINISHED"
)

var (
	Cfg *conf.Conf
)

func GetOsInfo(key string) string {
	var (
		output string
		err    error
	)
	cmd := "uname -" + key
	if output, err = RunCmd(cmd); err != nil {
		logrus.Errorf("get info by cmd[uname -v] failed.err=%v", err)
		return ""
	}
	return output
}

func RunCmd(cmd string) (string, error) {
	var (
		output []byte
		err    error
	)
	if output, err = exec.Command("/bin/sh", "-c", cmd).CombinedOutput(); err != nil {
		return "", err
	}
	return string(output), nil
}

func Fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ComposeResults(results *apistructs.TestResults) error {
	err := composeEnv(results)
	if err != nil {
		return err
	}

	if len(results.CommitID) > 6 {
		results.Name = fmt.Sprintf("ut-%s", results.CommitID[:6])
	} else {
		results.Name = fmt.Sprintf("ut-%s", results.CommitID)
	}

	results.Status = UTSatus
	results.Type = apistructs.UT

	return nil
}

func composeEnv(results *apistructs.TestResults) error {
	results.OperatorID = Cfg.OperatorID
	results.OperatorName = Cfg.OperatorName
	results.ApplicationID = int64(Cfg.AppID)
	results.ProjectID = int64(Cfg.ProjectID)
	results.ApplicationName = Cfg.AppName
	results.BuildID = Cfg.BuildID
	results.GitRepo = Cfg.GittarRepo
	results.Branch = Cfg.GittarBranch
	results.CommitID = Cfg.GittarCommit
	results.Workspace = Cfg.Workspace
	results.UUID = Cfg.UUID

	return nil
}

func ChangeWorkDir(codePath string) error {
	var (
		contextDir string
		err        error
	)
	contextDir = codePath

	if err = os.Chdir(contextDir); err != nil {
		return errors.Errorf("Change code_path: %s failed.", contextDir)
	}

	return nil
}

// ExecuteCmd 执行命令，并设置超时时间
func ExecuteCmd(cmdContent string) error {
	cmd := exec.Command("/bin/sh", "-c", cmdContent)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	finish := make(chan struct{})
	var err error
	go func() {
		err = cmd.Run()
		finish <- struct{}{}
	}()

	select {
	case <-time.After(180 * time.Minute):
		logrus.Warning("timeout to execute command(more than 180 minutes), exit")
		if err = cmd.Process.Kill(); err != nil {
			return errors.Errorf("failed to kill execute command process, cmd: %+v, (%+v)", cmd, err)
		}
	case <-finish:
		if err != nil {
			return errors.Errorf("failed to execute command, cmd: %+v, (%+v)", cmd, err)
		}
	}

	return nil
}

// ExecuteCmdOutput 执行命令，并设置超时时间, 返回结果数据
func ExecuteCmdOutput(cmdContent string) ([]byte, error) {
	var (
		err     error
		content []byte
	)

	cmd := exec.Command("/bin/sh", "-c", cmdContent)

	finish := make(chan struct{})
	go func() {
		content, err = cmd.CombinedOutput()
		finish <- struct{}{}
	}()

	select {
	case <-time.After(180 * time.Minute):
		logrus.Warning("timeout to execute command(more than 180 minutes), exit")
		if err = cmd.Process.Kill(); err != nil {
			return nil, errors.Errorf("failed to kill execute command process, cmd: %+v, (%+v)", cmd, err)
		}
	case <-finish:
		if err != nil {
			return nil, errors.Errorf("failed to execute command, cmd: %+v, (%+v)", cmd, err)
		}
	}

	return content, nil
}
