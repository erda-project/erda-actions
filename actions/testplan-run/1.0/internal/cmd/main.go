package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/testplan-run/1.0/internal/testplan-run"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("testplan-run Testing...")
	if err := testplan_run.Run(); err != nil {
		logrus.Errorf("testplan-run"+
			" failed, err: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
