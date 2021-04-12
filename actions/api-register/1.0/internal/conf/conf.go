package conf

// Conf action 入参
type Conf struct {
	WorkDir  string `env:"WORKDIR"`
	Metafile string `env:"METAFILE"`
	// Params

	SwaggerPath string `env:"ACTION_SWAGGER_PATH" required:"true"`
	ServiceName string `env:"ACTION_SERVICE_NAME" required:"true"`
	RuntimeID   string `env:"ACTION_RUNTIME_ID"`
	ServiceAddr string `env:"ACTION_SERVICE_ADDR"`

	// env
	OrgID             int64  `env:"DICE_ORG_ID" required:"true"`
	ClusterName       string `env:"DICE_CLUSTER_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID             int64  `env:"DICE_APPLICATION_ID" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE"`
	GittarBranch      string `env:"GITTAR_BRANCH"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
}
