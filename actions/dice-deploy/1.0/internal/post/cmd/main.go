package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice-deploy-service/1.0/dice-deploy-services"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("Canceling...")
	if err := dice_deploy_services.Cancel(); err != nil {
		logrus.Warning("Unable to cancel dice deploy, err: %v", err)
		return
	}
}
