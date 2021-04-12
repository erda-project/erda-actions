package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/api-register/1.0/internal/push"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	if err := push.Run(); err != nil {
		logrus.Errorf("Unable to register api to api-gateway, err: %+v", err)
		os.Exit(1)
	}
	logrus.Infof("api-register action success")

	os.Exit(0)
}
