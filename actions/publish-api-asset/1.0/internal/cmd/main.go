package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/publish-api-asset/1.0/internal/push"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	if err := push.Run(); err != nil {
		logrus.Errorf("publish api asset failed, err: %+v", err)
		os.Exit(1)
	}
	logrus.Infof("publish api asset action success")

	os.Exit(0)
}
