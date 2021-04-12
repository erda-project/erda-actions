package conf

import "github.com/erda-project/erda/pkg/envconf"

// Conf action 入参
type Conf struct {
	// wd & meta
	WorkDir  string `env:"WORKDIR"`
	MetaFile string `env:"METAFILE"`

	// env
	DiceOpenapiAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	DiceOpenapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	ProjectTestEnvID uint64 `env:"ACTION_PROJECT_TEST_ENV_ID"`
	ProjectID        uint64 `env:"ACTION_PROJECT_ID"`
	TestPlanID       uint64 `env:"ACTION_TEST_PLAN_ID"`
}

var (
	cfg Conf
)

func Load() error {
	return envconf.Load(&cfg)
}

func DiceOpenapiAddr() string {
	return cfg.DiceOpenapiAddr
}

func DiceOpenapiToken() string {
	return cfg.DiceOpenapiToken
}

func ProjectTestEnvID() uint64 {
	return cfg.ProjectTestEnvID
}

func ProjectID() uint64 {
	return cfg.ProjectID
}

func TestPlanID() uint64 {
	return cfg.TestPlanID
}
