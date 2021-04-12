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
	logrus.Info("Deploying services...")
	if err := dice_deploy_services.Run(); err != nil {
		logrus.Errorf("Unable to deploy application to dice, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
