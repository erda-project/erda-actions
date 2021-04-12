package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/conf"
)

func main() {

	logrus.SetOutput(os.Stdout)

	if err := conf.Initialize(); err != nil {
		bplog.Fatalf("failed to initialize conf, err: %v", err)
	}

	err := run.Execute()
	if err != nil {
		bplog.Fatal(err)
	}
}
