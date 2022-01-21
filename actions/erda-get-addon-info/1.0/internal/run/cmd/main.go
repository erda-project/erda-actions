package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/erda-get-addon-info/1.0/internal/erda-get-addon-info"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.Info("GetAddonConfigInfo...")
	if err := erda_get_addon_info.Run(); err != nil {
		logrus.Warning("Unable to get addon info, err: %v", err)
		os.Exit(1)
	}
}
