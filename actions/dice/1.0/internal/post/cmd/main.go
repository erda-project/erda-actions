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

	logrus.Info("Canceling...")
	if err := dice.Cancel(); err != nil {
		logrus.Warning("Unable to cancel dice deploy, err: %v", err)
		return
	}
}
