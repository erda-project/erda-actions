package main

import (
	"github.com/erda-project/erda-actions/actions/semgrep/1.0/internal/pkg"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.Init()

	cfg, err := pkg.Parse()
	if err != nil {
		logrus.Fatalf("failed to parse conf, err: %v\n", err)
	}

	semgrep, err := pkg.NewSemgrep(cfg)
	if err != nil {
		logrus.Fatalf("failed to initializing semgrep, err: %v\n", err)
	}

	if err := semgrep.Execute(); err != nil {
		logrus.Fatalf("failed to execute semgrep, err: %v\n", err)
	}
}
