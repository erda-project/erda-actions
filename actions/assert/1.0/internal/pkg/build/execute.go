package build

import (
	"fmt"

	"github.com/erda-project/erda-actions/actions/assert/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/assert"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/envconf"
)

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

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
		fmt.Printf("Assert Result:")
		fmt.Printf("  value: %v", jsonparse.JsonOneLine(v.Value))
		fmt.Printf("  assert: %v", v.Assert)
		fmt.Printf("  actualValue: %s", jsonparse.JsonOneLine(v.ActualValue))
		fmt.Printf("  success: %v", success)
		fmt.Printf("==========")
	}
	fmt.Printf("AllAssert Result: %v", allSuccess)
	if !allSuccess {
		return fmt.Errorf("asssert faild")
	}
	return nil
}
