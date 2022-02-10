package main

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/cancel"
)

func main() {
	log.Init()

	logrus.Printf("Canceling...")
	if err := cancel.Cancel(); err != nil {
		logrus.Warning("Unable to cancel dice deploy, err: %v", err)
		return
	}
}
