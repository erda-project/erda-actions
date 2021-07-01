package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
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

type SonarLogLevel string

const (
	SonarLogLevelINFO  SonarLogLevel = "INFO"
	SonarLogLevelDEBUG SonarLogLevel = "DEBUG"
	SonarLogLevelTRACE SonarLogLevel = "TRACE"
)

func (l SonarLogLevel) Valid() bool {
	switch l {
	case SonarLogLevelINFO, SonarLogLevelDEBUG, SonarLogLevelTRACE:
		return true
	default:
		return false
	}
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

	// used to invoke openapi
	OpenAPIAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	OpenAPIToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`

	DiceClusterName string `env:"DICE_CLUSTER_NAME" required:"true"`

	// metafile
	MetaFile string `env:"METAFILE"`

	LogID string `env:"TERMINUS_DEFINE_TAG"`

	// params
	ActionParams ActionParams
}

type ActionParams struct {
	// CodeDir 执行 sonar-scanner 的目录
	// +required
	CodeDir string `env:"ACTION_CODE_DIR" required:"true"`

	// language
	// go, java
	Language Language `env:"ACTION_LANGUAGE" requied:"true"`
	// java
	// +optional
	SonarJavaBinaries string `env:"ACTION_SONAR_JAVA_BINARIES"`

	// SonarHostURL sonar 服务器地址，用户可以手动指定。若不填写，则使用平台提供的 sonar 服务
	SonarHostURL string `env:"ACTION_SONAR_HOST_URL"`
	// SonarLogin is the login or authentication token of a SonarQube user with Execute Analysis permission on the project
	SonarLogin string `env:"ACTION_SONAR_LOGIN"`
	// SonarPassword is the password that goes with the sonar.login username. This should be left blank if an authentication token is being used
	SonarPassword string `env:"ACTION_SONAR_PASSWORD"`

	// sonar configs
	SonarExclusions string        `env:"ACTION_SONAR_EXCLUSIONS"`
	SonarLogLevel   SonarLogLevel `env:"SONAR_LOG_LEVEL" default:"INFO"`

	// Project
	// +optional
	ProjectKey string `env:"ACTION_PROJECT_KEY"`
	// DeleteProject after analysis finished
	// +optional
	DeleteProject bool `env:"ACTION_DELETE_PROJECT"`

	// 使用平台配置
	UsePlatformQualityGate bool `env:"ACTION_USE_PLATFORM_QUALITY_GATE"`

	// QualityGate
	QualityGate []QualityGateCondition `env:"ACTION_QUALITY_GATE"`
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

	if params.SonarHostURL == "" {
		// get platform sonar credential
		credential, err := getPlatformSonarCredential(&cfg)
		if err != nil {
			return nil, err
		}
		params.SonarHostURL = credential.Server
		params.SonarLogin = credential.Token
		// no need password, just use token, it's ok
	}

	cfg.ActionParams = params

	return &cfg, nil
}

func getPlatformSonarCredential(cfg *Conf) (*apistructs.SonarCredential, error) {
	var buffer bytes.Buffer
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(cfg.OpenAPIAddr).
		Path("/api/qa/actions/get-sonar-credential").
		Param("clusterName", cfg.DiceClusterName).
		Header("Authorization", cfg.OpenAPIToken).
		Do().Do().Body(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform sonar credential, err: %v", err)
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("failed to get platform sonar credential, statusCode: %d, respBody: %s", resp.StatusCode(), buffer.String())
	}
	var getResp apistructs.SonarCredentialGetResponse
	respBody := buffer.String()
	if err := json.Unmarshal([]byte(respBody), &getResp); err != nil {
		return nil, fmt.Errorf("failed to parse platform sonar credential, err: %v, json string: %s", err, respBody)
	}
	return getResp.Data, nil
}
