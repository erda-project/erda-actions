package main

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda/pkg/envconf"
)

type Language string

const (
	LanguageGo   Language = "go"
	LanguageJava Language = "java"
	LanguageJS   Language = "js"
)

func (l Language) Supported() bool {
	switch l {
	case LanguageGo, LanguageJava, LanguageJS:
		return true
	default:
		return false
	}
}

func (l Language) String() string {
	return string(l)
}

type Conf struct {
	Debug bool `env:"ACTION_DEBUG" default:"false"`

	ProjectID    uint64 `env:"DICE_PROJECT_ID"`
	ProjectName  string `env:"DICE_PROJECT_NAME"`
	AppID        uint64 `env:"DICE_APPLICATION_ID"`
	AppName      string `env:"DICE_APPLICATION_NAME"`
	Workspace    string `env:"DICE_WORKSPACE"`
	GittarRepo   string `env:"GITTAR_REPO"`
	GittarBranch string `env:"GITTAR_BRANCH"`
	GittarCommit string `env:"GITTAR_COMMIT"`
	OperatorID   string `env:"DICE_OPERATOR_ID"`
	PipelineID   int64  `env:"PIPELINE_ID"`

	OrgId  uint64 `env:"DICE_ORG_ID"`
	UserID string `env:"DICE_USER_ID"`

	DiceClusterName string `env:"DICE_CLUSTER_NAME" required:"true"`

	// metafile
	MetaFile string `env:"METAFILE"`

	LogID string `env:"TERMINUS_DEFINE_TAG"`

	// used to invoke openapi
	OpenAPIAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	OpenAPIToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`

	// params
	ActionParams ActionParams
}

type ActionParams struct {
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
	ProjectKey string `env:"ACTION_PROJECT_KEY"`

	MustGateStatusOK bool `env:"ACTION_MUST_GATE_STATUS_OK"`
}

func Parse() (*Conf, error) {
	// platform envs
	var cfg Conf
	err := envconf.Load(&cfg)
	if err != nil {
		return nil, err
	}

	// action params
	var params ActionParams
	err = envconf.Load(&params)
	if err != nil {
		return nil, err
	}

	cfg.ActionParams = params

	if cfg.ActionParams.SonarHostURL == "" {
		app, err := getApplication(&cfg)
		if err != nil {
			return nil, errors.Errorf("failed to get application: %v", err)
		}
		if app.SonarConfig == nil {
			return nil, errors.Errorf("application %s has no sonar config", app.Name)
		}
		cfg.ActionParams.SonarHostURL = app.SonarConfig.Host
		cfg.ActionParams.SonarLogin = app.SonarConfig.Token
		cfg.ActionParams.ProjectKey = app.SonarConfig.ProjectKey
	}

	return &cfg, nil
}
