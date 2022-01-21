package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/erda-add-app-config/1.0/internal/erda-add-app-config"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.Info("ReDeploying...")
	if err := erda_add_app_config.Run(); err != nil {
		logrus.Warning("Unable to add config to application, err: %v", err)
		os.Exit(1)
	}
}
