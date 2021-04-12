package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice-deploy-addon/1.0/internal/dice-deploy-addons"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.Info("Deploying addons...")
	if err := dice_deploy_addons.Run(); err != nil {
		logrus.Errorf("Unable to deploy application to dice, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
