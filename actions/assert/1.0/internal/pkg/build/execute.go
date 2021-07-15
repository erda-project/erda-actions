package build

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/assert/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/assert"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/envconf"
)

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)
	logrus.SetOutput(os.Stdout)
	if err := build(cfg); err != nil {
		return err
	}

	return nil
}

func build(cfg conf.Conf) error {
	var allSuccess = true
	for _, v := range cfg.Assert {
		success, err := assert.DoAssert(v.ActualValue, v.Assert, jsonparse.JsonOneLine(v.Value))
		if err != nil || !success {
			allSuccess = false
		}
		// to assert
		logrus.Infof("Assert Result:")
		logrus.Infof("  value: %v", jsonparse.JsonOneLine(v.Value))
		logrus.Infof("  assert: %v", v.Assert)
		logrus.Infof("  actualValue: %s", jsonparse.JsonOneLine(v.ActualValue))
		logrus.Infof("  success: %v", success)
		logrus.Infof("==========")
	}
	logrus.Infof("AllAssert Result: %v", allSuccess)
	if !allSuccess {
		return fmt.Errorf("asssert faild")
	}
	return nil
}
