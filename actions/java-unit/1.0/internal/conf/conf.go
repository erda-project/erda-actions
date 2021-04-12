package conf

// Conf action 入参
type Conf struct {
	// wd & meta
	WorkDir  string `env:"WORKDIR"`
	MetaFile string `env:"METAFILE"`

	// nexus
	NexusUrl      string `env:"BP_NEXUS_URL"`
	NexusUsername string `env:"BP_NEXUS_USERNAME"`
	NexusPassword string `env:"BP_NEXUS_PASSWORD"`

	// params
	Path    string `env:"ACTION_PATH"`
}
