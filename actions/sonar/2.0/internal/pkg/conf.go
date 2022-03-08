package pkg

import "github.com/erda-project/erda-actions/pkg/envconf"

type ActionParams struct {
	Debug bool `env:"ACTION_DEBUG" default:"false"`
	// CodeDir 执行 sonar-scanner 的目录
	// +required
	CodeDir string `env:"ACTION_CODE_DIR" required:"true"`

	// SonarHostURL sonar 服务器地址，用户可以手动指定。若不填写，则使用平台提供的 sonar 服务
	SonarHostURL string `env:"ACTION_SONAR_HOST_URL"`
	// SonarLogin is the login or authentication token of a SonarQube user with Execute Analysis permission on the project
	SonarLogin string `env:"ACTION_SONAR_LOGIN"`
	// SonarPassword is the password that goes with the sonar.login username. This should be left blank if an authentication token is being used
	SonarPassword string `env:"ACTION_SONAR_PASSWORD"`

	// Project
	// +optional
	ProjectKey string `env:"ACTION_SONAR_PROJECT_KEY"`

	MustGateStatusOK bool `env:"ACTION_MUST_GATE_STATUS_OK"`
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
