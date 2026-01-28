package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/mcp-register/1.0/internal/pkg/runner"
	"github.com/erda-project/erda-actions/pkg/log"
)

func main() {
	log.Init()

	logrus.Printf("Registing...")
	if err := runner.Run(); err != nil {
		logrus.Errorf("Unable to register mcp server, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
