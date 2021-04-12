package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice-deploy-redeploy/1.0/internal/dice-deploy-redeploy"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.Info("ReDeploying...")
	if err := dice_deploy_redeploy.Run(); err != nil {
		logrus.Warning("Unable to cancel dice deploy, err: %v", err)
		os.Exit(1)
	}
}
