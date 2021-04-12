package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/golang/1.0/internal/pkg/build"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	if err := build.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "Golang Action failed, err: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Golang Action success\n")
	os.Exit(0)
}
