package conf

type Conf struct {
	// 平台环境变量
	MetaFile  string `env:"METAFILE"`
	WorkDir   string `env:"WORKDIR"`
	UploadDir string `env:"UPLOADDIR"`

	NexusURL      string `env:"BP_NEXUS_URL" default:"https://repo.terminus.io"`
	NexusUsername string `env:"BP_NEXUS_USERNAME" default:"readonly"`
	NexusPassword string `env:"BP_NEXUS_PASSWORD" default:"Hello1234"`

	Memory float64 `env:"PIPELINE_LIMITED_MEM"`

	Debug bool `env:"DEBUG" default:"false"`

	// 用户指定
	CodeDir string `env:"ACTION_CODE_DIR" required:"true"`
}
