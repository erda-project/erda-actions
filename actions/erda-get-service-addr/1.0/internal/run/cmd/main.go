package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/erda-get-service-addr/1.0/internal/erda-get-service-addr"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.Info("Get Runtime Services Addr ...")
	if err := erda_get_service_addr.Run(); err != nil {
		logrus.Warning("Unable to get runtime services addr, err: %v", err)
		os.Exit(1)
	}
}
