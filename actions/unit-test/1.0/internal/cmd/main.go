package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/base"
	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/ut"
	"github.com/erda-project/erda-actions/pkg/detect/bptype"
	"github.com/erda-project/erda/pkg/envconf"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("Starting Unit Test...")

	ut := ut.NewUt()

	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		base.Fatal("conf not property", err)
	}
	base.Cfg = &cfg

	if err := bptype.RenderConfigToDir("/root/.m2"); err != nil {
		base.Fatal("render config", err)
	}
	if err := bptype.RenderConfigToDir("/root/.gradle"); err != nil {
		base.Fatal("render config", err)
	}

	if err := ut.UnitTest(); err != nil {
		base.Fatal("ut failed", err)
	}
}
