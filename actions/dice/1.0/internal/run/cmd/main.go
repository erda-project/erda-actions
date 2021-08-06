package main

import (
	"github.com/erda-project/erda-actions/pkg/log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice/1.0/internal/dice"
)

func main() {
	log.Init()

	logrus.Printf("Deploying...")
	if err := dice.Run(); err != nil {
		logrus.Errorf("Unable to deploy application to dice, err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
