package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/testscene-run/1.0/internal/testscene-run"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("testscene-run Testing...")
	if err := testscene_run.Run(); err != nil {
		logrus.Errorf("testscene-run"+
			" failed, err: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
