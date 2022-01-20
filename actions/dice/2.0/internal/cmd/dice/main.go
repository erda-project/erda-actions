package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/dice"
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
