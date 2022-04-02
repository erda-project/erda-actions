package pkg

import (
	"github.com/erda-project/erda-actions/pkg/meta"
)

type MultiMerge struct {
	cfg     *Conf
	results *meta.ResultMetaCollector
}

func NewMultiMerge(cfg *Conf) (*MultiMerge, error) {
	m := &MultiMerge{cfg: cfg, results: meta.NewResultMetaCollector()}
	return m, nil
}
