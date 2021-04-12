package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func produceCoverageArgs(lan, goDir, contextDir string) (string, error) {
	switch lan {
	case "js":
		if err := produceJsCoverageFile(); err != nil {
			return "", err
		}

		return " -Dsonar.javascript.lcov.reportPaths='coverage/lcov.info' ", nil
	case "golang":
		if err := produceGolangCoverageFile(goDir, contextDir); err != nil {
			return "", err
		}

		return " -Dsonar.language=go -Dsonar.exclusions='**/vendor/**,config/**,docs/**,**/*_test.go' -Dsonar.tests='.' -Dsonar.test.inclusions='**/*_test.go' -Dsonar.test.exclusions='**/vendor/**,config/**,docs/**' -Dsonar.go.coverage.reportPaths='coverage.out' ", nil
	}

	return "", nil
}

func produceGolangCoverageFile(goDir, contextDir string) error {
	var (
		err error
	)
	if goDir == "" {
		return errors.New("need go_dir")
	}

	goWorkSpace := filepath.Join("/go/src", goDir)

	// make go dir and cp
	if err = executeCmd(fmt.Sprintf("mkdir -p %s; cp -rf %s %s",
		filepath.Dir(goWorkSpace), contextDir, goWorkSpace)); err != nil {
		return err
	}

	if err = os.Chdir(goWorkSpace); err != nil {
		return errors.Wrapf(err, "failed to change directory, path: %s", goWorkSpace)
	}

	if err = executeCmd("go test ./... -coverprofile=coverage.out"); err != nil {
		return err
	}

	_, err = os.Stat("coverage.out")
	if os.IsNotExist(err) {
		return errors.New("not exist coverage.out")
	}

	return nil
}

func produceJsCoverageFile() error {
	// install npm mocha, istanbul
	if err := executeCmd("npm install --global mocha istanbul; npm install"); err != nil {
		return err
	}

	// start test.
	if err := executeCmd("istanbul cover _mocha *test*"); err != nil {
		if err = executeCmd("npm run coverage"); err != nil {
			return err
		}
	}

	_, err := os.Stat("coverage/lcov.info")
	if os.IsNotExist(err) {
		return errors.New("not exist coverage/lcov.info")
	}

	return nil
}

// ExecuteCmd 执行命令，并设置超时时间
func executeCmd(cmdContent string) error {
	cmd := exec.Command("/bin/sh", "-c", cmdContent)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr

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
