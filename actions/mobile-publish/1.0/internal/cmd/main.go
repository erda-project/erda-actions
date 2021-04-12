package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/mobile-publish/1.0/internal/push"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	if err := push.Run(); err != nil {
		logrus.Errorf("Unable to push mobile to publish-item, err: %v", err)
		os.Exit(1)
	}
	logrus.Infof("mobile-publish action success")

	os.Exit(0)
}
