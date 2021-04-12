package conf

import (
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/strutil"
)

// Conf action 入参
type Conf struct {
	WorkDir  string `env:"WORKDIR"`
	Metafile string `env:"METAFILE"`
	// Params
	UsecaseID uint64 `env:"ACTION_USECASE_ID"`
	APIID     uint64 `env:"ACTION_API_ID"`
	APIIDs    string `env:"ACTION_API_IDS"`
	// env
	DiceOpenapiAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	DiceOpenapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	ProjectTestEnvID uint64 `env:"PROJECT_TEST_ENV_ID"`
}

var (
	cfg Conf
)

func Load() error {
	return envconf.Load(&cfg)
}

func APIIDs() []uint64 {
	if cfg.APIID != 0 {
		return []uint64{cfg.APIID}
	}
	if cfg.APIIDs == "" {
		return nil
	}
	v := strutil.Split(cfg.APIIDs, ",", true)
	ids := make([]uint64, 0, len(v))
	for _, i := range v {
		id, err := strutil.Atoi64(i)
		if err != nil {
			continue
		}
		ids = append(ids, uint64(id))
	}
	return ids
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

func UsecaseID() uint64 {
	return cfg.UsecaseID
}
