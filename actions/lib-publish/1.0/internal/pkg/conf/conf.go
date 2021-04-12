package conf

// Conf action 入参
type Conf struct {
	WorkDir  string `env:"WORKDIR"`
	Metafile string `env:"METAFILE"`
	// Params

	Context string `env:"ACTION_WORKDIR" required:"true"` // spec.yml & README.md 所在目录

	// env
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE"`
	UserID            string `env:"DICE_USER_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
}
