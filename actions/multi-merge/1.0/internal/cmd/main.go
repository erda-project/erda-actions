package main

import (
	"github.com/erda-project/erda-actions/actions/multi-merge/1.0/internal/pkg"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.Init()

	cfg, err := pkg.Parse()
	if err != nil {
		logrus.Fatalf("failed to parse conf, err: %v\n", err)
	}

	merge, err := pkg.NewMultiMerge(cfg)
	if err != nil {
		logrus.Fatalf("failed to create multi-merge, err: %v\n", err)
	}
	if err := merge.Execute(); err != nil {
		logrus.Fatalf("failed to execute multi-merge, err: %v\n", err)
	}
}
