package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/app-run/1.0/internal/appCancel/cancel"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("Canceling...")
	if err := cancel.Cancel(); err != nil {
		logrus.Warning("Unable to cancel app-run deploy, err: %v", err)
		return
	}
}
