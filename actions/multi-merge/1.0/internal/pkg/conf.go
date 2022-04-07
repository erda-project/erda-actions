package pkg

import (
	"github.com/erda-project/erda-actions/pkg/envconf"
	"github.com/pkg/errors"
)

type Repo struct {
	Uri      string   `json:"uri"`
	Branches []string `json:"branches"`
}

type GitConfig struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ActionParams struct {
	DestRepo   string      `env:"ACTION_DEST_REPO"`
	DestBranch string      `env:"ACTION_DEST_BRANCH"`
	Username   string      `env:"ACTION_USERNAME"`
	Password   string      `env:"ACTION_PASSWORD"`
	Repos      []Repo      `env:"ACTION_REPOS"`
	GitConfigs []GitConfig `env:"ACTION_GIT_CONFIG"`
}

type Conf struct {
	PlatformParams envconf.PlatformParams
	ActionParams
}

func Parse() (*Conf, error) {
	conf := &Conf{}
	platform, err := envconf.NewPlatformParams()
	if err != nil {
		return nil, err
	}
	conf.PlatformParams = platform
	actionParam := ActionParams{}
	if err := envconf.Load(&actionParam); err != nil {
		return nil, err
	}
	conf.ActionParams = actionParam
	if conf.DestRepo == "" {
		conf.DestRepo = conf.PlatformParams.GittarRepo
	}
	if conf.DestRepo == "" {
		return nil, errors.Errorf("destination repo is not set")
	}
	if conf.DestBranch == "" {
		conf.DestBranch = conf.PlatformParams.GittarBranch
	}
	if conf.DestBranch == "" {
		return nil, errors.Errorf("destination branch is not set")
	}
	return conf, nil
}
