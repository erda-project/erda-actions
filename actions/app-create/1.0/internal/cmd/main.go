package main

import (
	"os"

	"github.com/erda-project/erda-actions/actions/app-create/1.0/internal/appCreate"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("app-create Testing...")

	if err := appCreate.Run(); err != nil {
		logrus.Errorf("app-create failed, err: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
