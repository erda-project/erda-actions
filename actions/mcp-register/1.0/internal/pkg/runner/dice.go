package runner

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/mcp-register/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/mcp-register/1.0/internal/pkg/register"
)

func Run() error {
	// parse config
	cfg, err := conf.HandleConf()
	if err != nil {
		logrus.Errorf("failed to handle conf, err: %v", err)
		return err
	}

	r := register.New(&cfg)
	err = r.Register()
	if err != nil {
		return err
	}

	return nil
}
