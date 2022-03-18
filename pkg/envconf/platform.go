package envconf

type PlatformParams struct {
	ProjectID    uint64 `env:"DICE_PROJECT_ID"`
	ProjectName  string `env:"DICE_PROJECT_NAME"`
	AppID        uint64 `env:"DICE_APPLICATION_ID"`
	AppName      string `env:"DICE_APPLICATION_NAME"`
	Workspace    string `env:"DICE_WORKSPACE"`
	GittarRepo   string `env:"GITTAR_REPO"`
	GittarBranch string `env:"GITTAR_BRANCH"`
	GittarCommit string `env:"GITTAR_COMMIT"`
	OperatorID   string `env:"DICE_OPERATOR_ID"`
	PipelineID   int64  `env:"PIPELINE_ID"`

	OrgID  uint64 `env:"DICE_ORG_ID"`
	UserID string `env:"DICE_USER_ID"`

	DiceClusterName string `env:"DICE_CLUSTER_NAME" required:"true"`

	// metafile
	MetaFile string `env:"METAFILE"`

	LogID string `env:"TERMINUS_DEFINE_TAG"`

	// used to invoke openapi
	OpenAPIAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	OpenAPIToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
}

func NewPlatformParams() (PlatformParams, error) {
	platform := PlatformParams{}
	err := Load(&platform)
	if err != nil {
		return platform, err
	}
	return platform, nil
}
