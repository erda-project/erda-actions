package conf

import (
	"github.com/erda-project/erda/pkg/envconf"
)

// Conf action 入参
type Conf struct {
	ProcessorId         []string `env:"ACTION_PROCESSOR" required:"true"`
	WaitTimeIntervalSec int      `env:"ACTION_WAIT_TIME_INTERVAL_SEC" default:"5"`
	// env
	MetaFile         string `env:"METAFILE"`
	WorkDir          string `env:"WORKDIR" default:"."`
	DiceVersion      string `env:"DICE_VERSION"`
	PipelineId       uint64 `env:"PIPELINE_ID"`
	OrgId            uint64 `env:"DICE_ORG_ID"`
	TaskId           uint64 `env:"PIPELINE_TASK_ID"`
	ApplicationId    uint64 `env:"DICE_APPLICATION_ID"`
	ProjectId        uint64 `env:"DICE_PROJECT_ID"`
	SponsorId        string `env:"DICE_USER_ID"`
	CommitId         string `env:"GITTAR_COMMIT"`
	ProjectName      string `env:"DICE_PROJECT_NAME"`
	ApplicationName  string `env:"DICE_APPLICATION_NAME"`
	BranchName       string `env:"GITTAR_BRANCH"`
	DiceOpenapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
}

var (
	cfg Conf
)

func Load() error {
	return envconf.Load(&cfg)
}
func PipelineId() uint64 {
	return cfg.PipelineId
}
func OrgId() uint64 {
	return cfg.OrgId
}
func TaskId() uint64 {
	return cfg.TaskId
}

func ApplicationId() uint64 {
	return cfg.ApplicationId
}
func ProjectId() uint64 {
	return cfg.ProjectId
}
func SponsorId() string {
	return cfg.SponsorId
}
func CommitId() string {
	return cfg.CommitId
}
func ProjectName() string {
	return cfg.ProjectName
}
func ApplicationName() string {
	return cfg.ApplicationName
}
func BranchName() string {
	return cfg.BranchName
}
func DiceOpenapiToken() string {
	return cfg.DiceOpenapiToken
}

func MetaFile() string {
	return cfg.MetaFile
}

func DiceOpenapiAddr() string {
	return cfg.DiceOpenapiAddr
}

func ProcessorId() []string {
	return cfg.ProcessorId
}

func WaitTimeIntervalSec() int {
	return cfg.WaitTimeIntervalSec
}
