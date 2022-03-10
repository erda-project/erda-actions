package config

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/pkg/envconf"
)

var conf *Config

type Config struct {
	MetaFile string `env:"METAFILE"`

	// basic envs
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	UserID            string `env:"DICE_USER_ID"`

	// action parameters
	Artifacts   []Artifact `env:"ACTION_ARTIFACTS" required:"true"`
	WaitMinutes int        `env:"ACTION_MINUTES" default:"1"`
}

type Artifact struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

func New() (*Config, error) {
	if conf != nil {
		return conf, nil
	}

	conf = new(Config)
	if err := envconf.Load(conf); err != nil {
		return nil, errors.Wrap(err, "failed to Load envs")
	}

	return conf, nil
}
