package conf

// Conf js action param collection
type Conf struct {
	MetaFile string `env:"METAFILE"`
	WorkDir  string `env:"WORKDIR"`
	// 用户指定
	Context     string `env:"ACTION_WORKDIR" required:"true"` // npm publish 目录
	NpmRegistry string `env:"ACTION_REGISTRY" required:"true"`
	NpmUsername string `env:"ACTION_USERNAME" required:"true"`
	NpmPassword string `env:"ACTION_PASSWORD" required:"true"`
}
