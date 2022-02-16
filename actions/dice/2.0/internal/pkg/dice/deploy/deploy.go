package deploy

import (
	"sync"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/dice/store"
)

type Deploy interface {
	Do() (string, map[string]*common.DeployResult, error)
	StatusCheck(orderId string, result map[string]*common.DeployResult, timeOut int) error
}

type Option func(d *deploy)
type deploy struct {
	cfg        *conf.Conf
	store      store.Store
	statusLock sync.Mutex
}

func New(opts ...Option) Deploy {
	d := deploy{}
	for _, opt := range opts {
		opt(&d)
	}
	return &d
}

func WithConf(c *conf.Conf) Option {
	return func(d *deploy) {
		d.cfg = c
	}
}

func WithStore(s store.Store) Option {
	return func(d *deploy) {
		d.store = s
	}
}
