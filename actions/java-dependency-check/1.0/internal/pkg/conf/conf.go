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

	// 用户指定
	CodeDir              string `env:"ACTION_CODE_DIR" required:"true"`
	Debug                bool   `env:"ACTION_DEBUG" default:"false"`
	AutoUpdateNVD        bool   `env:"ACTION_AUTO_UPDATE_NVD" default:"false"`
	MavenPluginVersion   string `env:"ACTION_MAVEN_PLUGIN_VERSION" default:"6.3.1"`
	MavenSettingsXMLPath string `env:"ACTION_MAVEN_SETTINGS_XML_PATH" required:"false"`
}
