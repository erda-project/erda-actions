package pkg

import "github.com/erda-project/erda-actions/pkg/envconf"

type ActionParams struct {
	Debug bool `env:"ACTION_DEBUG" default:"false"`
	// CodeDir 执行 semgrep ci 的目录
	// +required
	CodeDir string `env:"ACTION_CODE_DIR" required:"true"`

	Config string `env:"ACTION_CONFIG" required:"true"`

	Format string `env:"ACTION_FORMAT"`

	Args []string `env:"ACTION_ARGS"`
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
	return conf, nil
}
