package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/manual-review/1.0/internal/manualReview"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("manualReview Testing...")
	if err := manualReview.Run(); err != nil {
		logrus.Errorf("manualReview failed, err: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
