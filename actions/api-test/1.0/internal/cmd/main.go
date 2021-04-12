package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/api-test/1.0/internal/apitest"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("API Testing...")
	if err := apitest.Run(); err != nil {
		logrus.Errorf("API Test failed, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
