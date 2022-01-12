package dice

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda/pkg/envconf"
)

type conf struct {
	OrgID        uint64 `env:"DICE_ORG_ID"`
	ProjectID    uint64 `env:"DICE_PROJECT_ID"`
	AppID        uint64 `env:"DICE_APPLICATION_ID"`
	Workspace    string `env:"DICE_WORKSPACE"`
	GittarBranch string `env:"GITTAR_BRANCH"`
	ClusterName  string `env:"DICE_CLUSTER_NAME"`
	OperatorID   string `env:"DICE_OPERATOR_ID"`

	// used to invoke openapi
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR"`
	DiceOpenapiToken  string `env:"DICE_OPENAPI_TOKEN"`
	InternalClient    string `env:"DICE_INTERNAL_CLIENT"`
	UserID            string `env:"DICE_USER_ID"`

	// wd & meta
	WorkDir  string `env:"WORKDIR"`
	MetaFile string `env:"METAFILE"`

	PipelineBuildID uint64 `env:"PIPELINE_ID"`
	PipelineTaskID  uint64 `env:"PIPELINE_TASK_ID"`

	// params
	ReleaseID         string `env:"ACTION_RELEASE_ID"`
	ReleaseIDPath     string `env:"ACTION_RELEASE_ID_PATH"`
	TimeOut           int    `env:"ACTION_TIME_OUT"`
	Callback          string `env:"ACTION_CALLBACK"`
	EdgeLocation      string `env:"ACTION_EDGE_LOCATION"`
	AssignedWorkspace string `env:"ACTION_WORKSPACE"`
}

// HiddenActionParams value passed from user, but not defined in spec.yml
type HiddenActionParams struct {
	OrgID        uint64 `env:"ACTION_DICE_ORG_ID"`
	ProjectID    uint64 `env:"ACTION_DICE_PROJECT_ID"`
	AppID        uint64 `env:"ACTION_DICE_APPLICATION_ID"`
	Workspace    string `env:"ACTION_DICE_WORKSPACE"`
	GittarBranch string `env:"ACTION_GITTAR_BRANCH"`
	ClusterName  string `env:"ACTION_DICE_CLUSTER_NAME"`
}

func HandleConf() (conf, error) {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return conf{}, err
	}
	var hiddenActionParams HiddenActionParams
	if err := envconf.Load(&hiddenActionParams); err != nil {
		return conf{}, err
	}

	// assign action params if not empty
	if hiddenActionParams.OrgID > 0 {
		cfg.OrgID = hiddenActionParams.OrgID
	}
	if hiddenActionParams.ProjectID > 0 {
		cfg.ProjectID = hiddenActionParams.OrgID
	}
	if hiddenActionParams.AppID > 0 {
		cfg.AppID = hiddenActionParams.AppID
	}
	if hiddenActionParams.Workspace != "" {
		cfg.Workspace = hiddenActionParams.Workspace
	}
	if hiddenActionParams.GittarBranch != "" {
		cfg.GittarBranch = hiddenActionParams.GittarBranch
	}
	if hiddenActionParams.ClusterName != "" {
		cfg.ClusterName = hiddenActionParams.ClusterName
	}

	cfg.print()
	return cfg, nil
}

func (cfg *conf) print() {
	log.AddNewLine(1)
	logrus.Infof("config: ")
	logrus.Infof(" appID: %d", cfg.AppID)
	logrus.Infof(" clusterName: %s", cfg.ClusterName)
	logrus.Infof(" workspace: %s", cfg.Workspace)
	logrus.Infof(" gittarBranch: %s", cfg.GittarBranch)
	logrus.Infof(" operatorID: %s", cfg.OperatorID)
	log.AddLineDelimiter(" ")
}
