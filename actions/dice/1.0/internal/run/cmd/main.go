package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice/1.0/internal/dice"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("Deploying...")
	if err := dice.Run(); err != nil {
		logrus.Errorf("Unable to deploy application to dice, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
