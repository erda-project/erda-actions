package main

import (
	"github.com/erda-project/erda-actions/actions/sonar/2.0/internal/pkg"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.Init()

	cfg, err := pkg.Parse()
	if err != nil {
		logrus.Fatalf("failed to parse conf, err: %v\n", err)
	}

	sonar, err := pkg.NewSonar(cfg)
	if err != nil {
		logrus.Fatalf("failed to initializing sonar, err: %v\n", err)
	}

	if err := sonar.Execute(); err != nil {
		logrus.Fatalf("failed to execute sonar, err: %v\n", err)
	}
}
