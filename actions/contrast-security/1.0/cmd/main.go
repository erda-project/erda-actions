package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/contrast-security/1.0/internal/pkg/build"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	logrus.Info("Contrast Security start...")
	if err := build.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "Contrast Secutiry execute failed, err: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
