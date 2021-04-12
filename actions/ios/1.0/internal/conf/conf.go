package conf

type Conf struct {
	MetaFile        string `env:"METAFILE"`
	WorkDir         string `env:"WORKDIR"`
	PipelineContext string `env:"CONTEXTDIR"`

	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	PipelineTaskID    string `env:"PIPELINE_TASK_ID" required:"true"`
	PipelineID        string `env:"PIPELINE_ID" required:"true"`

	Commands        []string             `env:"ACTION_COMMANDS"`
	Targets         []string             `env:"ACTION_TARGETS"`
	Context         string               `env:"ACTION_CONTEXT"`
	P12Cert         *P12CertFile         `env:"ACTION_P12_CERT"`
	MobileProvision *MobileProvisionFile `env:"ACTION_MOBILE_PROVISION"`

	// pipeline注入，镜像生成需要
	PipelineTaskLogID string `env:"PIPELINE_TASK_LOG_ID" `
}

type MobileProvisionFile struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

type P12CertFile struct {
	Source   string      `json:"source"`
	Password interface{} `json:"password"`
	Dest     string      `json:"dest"`
}
