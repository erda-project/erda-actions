package testplan_run

import (
	"github.com/erda-project/erda-actions/actions/testplan-run/1.0/internal/conf"
)

func Run() error {
	if err := conf.Load(); err != nil {
		return err
	}
	return handleAPIs()
}
