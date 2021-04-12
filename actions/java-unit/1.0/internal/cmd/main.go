package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/java-unit/1.0/internal/conf"
	"github.com/erda-project/erda-actions/pkg/detect/bptype"
	"github.com/erda-project/erda/pkg/envconf"
)

func main() {
	err := execute()
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}
}

func execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	if err := bptype.RenderConfigToDir("/root/.m2"); err != nil {
		logrus.Errorf("failed to render config, (%+v)", err)
		return err
	}

	cmdStr := "MAVEN_OPTS=-Xmx2016m mvn clean test -Dmaven.test.failure.ignore=true"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	if cfg.Path != "" && cfg.Path != "." {
		cmd.Dir = cfg.Path
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
