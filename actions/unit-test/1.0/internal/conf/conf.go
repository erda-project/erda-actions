package conf

type Conf struct {
	ProjectID          uint64 `env:"DICE_PROJECT_ID"`
	AppID              uint64 `env:"DICE_APPLICATION_ID"`
	AppName            string `env:"DICE_APPLICATION_NAME"`
	Workspace          string `env:"DICE_WORKSPACE"`
	GittarRepo         string `env:"GITTAR_REPO"`
	GittarBranch       string `env:"GITTAR_BRANCH"`
	GittarCommit       string `env:"GITTAR_COMMIT"`
	ClusterName        string `env:"DICE_CLUSTER_NAME"`
	OperatorID         string `env:"DICE_OPERATOR_ID"`
	OperatorName       string `env:"DICE_OPERATOR_NAME"`
	BuildID            int64  `env:"PIPELINE_ID"`
	PipelineLimitedMem string `env:"PIPELINE_LIMITED_MEM" default:"1024"`
	UUID               string `env:"TERMINUS_DEFINE_TAG"`

	// used to invoke openapi
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR"`
	DiceOpenapiToken  string `env:"DICE_OPENAPI_TOKEN"`

	// wd & meta
	WorkDir  string `env:"WORKDIR"`
	MetaFile string `env:"METAFILE"`

	// nexus
	NexusUrl      string `env:"BP_NEXUS_URL"`
	NexusUsername string `env:"BP_NEXUS_USERNAME"`
	NexusPassword string `env:"BP_NEXUS_PASSWORD"`

	// params
	Context string `env:"ACTION_CONTEXT"`
	GoDir   string `env:"ACTION_GO_DIR"`
	Command string `env:"ACTION_COMMAND"`
}
