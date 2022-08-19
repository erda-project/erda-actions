package pkg

import (
	"github.com/erda-project/erda-actions/pkg/command"
	"github.com/erda-project/erda-actions/pkg/meta"
)

type Semgrep struct {
	cfg     *Conf
	cmd     *command.Cmd
	results *meta.ResultMetaCollector
}

func NewSemgrep(cfg *Conf) (*Semgrep, error) {
	sonar := Semgrep{cfg: cfg, cmd: command.NewCmd("semgrep"), results: meta.NewResultMetaCollector()}
	return &sonar, nil
}
