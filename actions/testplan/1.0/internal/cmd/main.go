package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/testplan/1.0/internal/testplan"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("Start exec test plan...")
	if err := testplan.Run(); err != nil {
		logrus.Errorf("failed to exec test plan, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
