package conf

// Conf action 入参
type Conf struct {
	// Params
	DisplayName string `env:"ACTION_DISPLAY_NAME"`
	AssetID     string `env:"ACTION_ASSET_ID" required:"true"`
	Version     string `env:"ACTION_VERSION"`
	SpecPath    string `env:"ACTION_SPEC_PATH" required:"true"`
	RuntimeID   uint64 `env:"ACTION_RUNTIME_ID" required:"true"`
	ServiceName string `env:"ACTION_SERVICE_NAME" required:"true"`

	// env
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	ProjectID         uint64 `env:"DICE_PROJECT_ID" required:"true"`
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	UserID            string `env:"DICE_USER_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
}
