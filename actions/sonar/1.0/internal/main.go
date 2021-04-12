package main

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
)

func main() {
	log.Init()

	cfg, err := Parse()
	if err != nil {
		logrus.Fatalf("failed to parse conf, err: %v\n", err)
	}

	SONAR := NewSonar(cfg.ActionParams.SonarHostURL, cfg.ActionParams.SonarLogin, cfg.ActionParams.SonarPassword)
	command := NewCommand(SONAR)

	if err := command.Analysis(cfg); err != nil {
		logrus.Fatalf("Sonar Analysis failed, err: %v", err)
	}
}
