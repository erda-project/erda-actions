package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/app-run/1.0/internal/appRun"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("app-run Testing...")
	if err := appRun.Run(); err != nil {
		logrus.Errorf("app-create failed, err: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
