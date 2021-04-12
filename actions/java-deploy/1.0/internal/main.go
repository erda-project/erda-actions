package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/dlog"
)

func main() {
	logrus.SetOutput(os.Stdout)

	err := conf.LoadEnvConfig()
	if err != nil {
		dlog.Fatalf("failed to load env config, err: %v\n", err)
	}

	err = run.Execute()
	if err != nil {
		dlog.Fatal(err)
	}
}
