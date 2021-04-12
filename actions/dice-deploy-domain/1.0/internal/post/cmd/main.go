package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice-deploy-domain/1.0/internal/dice-deploy-domain"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("Canceling...")
	if err := dice_deploy_domains.Cancel(); err != nil {
		logrus.Warning("Unable to cancel dice deploy, err: %v", err)
		return
	}
}
