package main

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice/1.0/internal/dice"
	"github.com/erda-project/erda-actions/pkg/log"
)

func main() {
	log.Init()

	logrus.Printf("Canceling...")
	if err := dice.Cancel(); err != nil {
		logrus.Warning("Unable to cancel dice deploy, err: %v", err)
		return
	}
}
