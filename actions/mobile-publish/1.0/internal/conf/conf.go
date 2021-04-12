package conf

// Conf action 入参
type Conf struct {
	WorkDir  string `env:"WORKDIR"`
	Metafile string `env:"METAFILE"`
	// Params

	ReleaseID  string `env:"ACTION_RELEASE_ID"`
	ReadmeFile string `env:"ACTION_README_FILE"`
	SpecFile   string `env:"ACTION_SPEC_FILE"`

	// env
	UserID            string `env:"DICE_USER_ID" required:"true"`
	AppID             int64  `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
}
