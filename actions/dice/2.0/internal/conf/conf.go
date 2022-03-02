package conf

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda/pkg/envconf"
)

type Conf struct {
	OrgID        uint64 `env:"DICE_ORG_ID"`
	ProjectID    uint64 `env:"DICE_PROJECT_ID"`
	AppID        uint64 `env:"DICE_APPLICATION_ID"`
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
	ReleaseID           string `env:"ACTION_RELEASE_ID"`
	ReleaseIDPath       string `env:"ACTION_RELEASE_ID_PATH"`
	ReleaseName         string `env:"ACTION_RELEASE_NAME"`
	ReleaseTye          string `env:"ACTION_TYPE"`
	ApplicationName     string `env:"ACTION_APPLICATION_NAME"`
	DeployWithoutBranch bool   `env:"ACTION_DEPLOY_WITHOUT_BRANCH"`
	TimeOut             int    `env:"ACTION_TIME_OUT"`
	Callback            string `env:"ACTION_CALLBACK"`
	EdgeLocation        string `env:"ACTION_EDGE_LOCATION"`
	AssignedWorkspace   string `env:"ACTION_WORKSPACE"`
}

// HiddenActionParams value passed from user, but not defined in spec.yml
type HiddenActionParams struct {
	OrgID        uint64 `env:"ACTION_DICE_ORG_ID"`
	ProjectID    uint64 `env:"ACTION_DICE_PROJECT_ID"`
	AppID        uint64 `env:"ACTION_DICE_APPLICATION_ID"`
	GittarBranch string `env:"ACTION_GITTAR_BRANCH"`
	ClusterName  string `env:"ACTION_DICE_CLUSTER_NAME"`
}

func HandleConf() (Conf, error) {
	var cfg Conf
	if err := envconf.Load(&cfg); err != nil {
		return Conf{}, err
	}
	var hiddenActionParams HiddenActionParams
	if err := envconf.Load(&hiddenActionParams); err != nil {
		return Conf{}, err
	}

	// assign action params if not empty
	if hiddenActionParams.OrgID > 0 {
		cfg.OrgID = hiddenActionParams.OrgID
	}
	if hiddenActionParams.ProjectID > 0 {
		cfg.ProjectID = hiddenActionParams.ProjectID
	}
	if hiddenActionParams.AppID > 0 {
		cfg.AppID = hiddenActionParams.AppID
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

func (cfg *Conf) print() {
	log.AddNewLine(1)
	logrus.Infof("config: ")
	logrus.Infof(" appID: %d", cfg.AppID)
	logrus.Infof(" projectId: %d", cfg.ProjectID)
	logrus.Infof(" clusterName: %s", cfg.ClusterName)
	logrus.Infof(" gittarBranch: %s", cfg.GittarBranch)
	logrus.Infof(" operatorID: %s", cfg.OperatorID)
	logrus.Infof(" deployWithoutBranch: %v", cfg.DeployWithoutBranch)
	if cfg.AssignedWorkspace != "" {
		logrus.Infof(" assignedWorkspace: %v", cfg.AssignedWorkspace)
	}
	log.AddLineDelimiter(" ")
}
