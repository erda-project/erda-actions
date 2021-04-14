package config

import (
	"strconv"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
)

// out params keys
const (
	Commit  = "commit"
	MrID    = "mr_id"
	MrUrl   = "mr_url"
	Success = "success"
	Err     = "err"
	Warn    = "warn"
	Step    = "step"
)

const (
	DiceYmlPathFromSrcRepo             = "dice.yml"
	DiceYmlPathFromDstRepoVersionDir   = "releases/dice/dice.yml"
	MigrationPathFromSrcRepo           = ".dice/migrations"
	MigrationPathFromDstRepoVersionDir = "sqls"
)

var c *config

func init() {
	initLog()
}

type config struct {
	// 基本环境参数
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE" required:"true"`

	// 流水线配置
	PipelineDebugMode bool   `env:"PIPELINE_DEBUG_MODE"`
	PipelineID        string `env:"PIPELINE_ID"`
	PipelineTaskLogID string `env:"PIPELINE_TASK_LOG_ID"`
	PipelineTaskID    string `env:"PIPELINE_TASK_ID"`

	// Action 入参
	Workdir     string               `env:"ACTION_WORKDIR"`
	Dst         RepoInfo             `env:"ACTION_DST"`
	MRProcessor uint64               `env:"ACTION_MR_PROCESSOR"`
	Registry    *RegistryReplacement `env:"ACTION_REGISTRY_REPLACEMENT"`

	// 其他参数
	MetaFilename string `env:"METAFILE"`
}

type RepoInfo struct {
	ApplicationName string `json:"applicationName"`
	Branch          string `json:"branch"`
}

type RegistryReplacement struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func configuration() *config {
	if c == nil {
		c = new(config)
		if err := envconf.Load(c); err != nil {
			logrus.Fatalf("failed to load configuration, err: %v", err)
		}
		if c.Dst.Branch == "" {
			c.Dst.Branch = "master"
		}
	}

	return c
}

func OrdID() uint64 {
	return configuration().OrgID
}

func OpenapiToken() string {
	return configuration().CiOpenapiToken
}

func OpenapiPrefix() string {
	return configuration().DiceOpenapiPrefix
}

func Workdir() string {
	return configuration().Workdir
}

func DstRepo() RepoInfo {
	return configuration().Dst
}

func Metafile() string {
	return configuration().MetaFilename
}

func PipelineID() string {
	return configuration().PipelineID
}

func ProjectID() int64 {
	return configuration().ProjectID
}

func ProjectName() string {
	return configuration().ProjectName
}

func ApplicationID() uint64 {
	return configuration().AppID
}

func ApplicationName() string {
	return configuration().AppName
}

func DstApplicationName() string {
	if configuration().Dst.ApplicationName == "" {
		return "version"
	}
	return configuration().Dst.ApplicationName
}

func DstRepoRefBranch() string {
	return configuration().Dst.Branch
}

func DstRepoBranch() string {
	return "feature/pipeline-" + configuration().PipelineID + "-" + configuration().PipelineTaskID
}

func MRProcessor() string {
	return strconv.FormatUint(configuration().MRProcessor, 10)
}

func Replacement() *RegistryReplacement {
	return configuration().Registry
}

func initLog() {
	log.Init()
	if configuration().PipelineDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
