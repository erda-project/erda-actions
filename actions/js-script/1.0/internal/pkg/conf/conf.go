package conf

// Conf js action param collection
type Conf struct {
	MetaFile string `env:"METAFILE"`
	WorkDir  string `env:"WORKDIR"`

	// 用户指定
	Context  string   `env:"ACTION_WORKDIR" required:"true"`
	Commands []string `env:"ACTION_COMMANDS" required:"true"`
	Targets  []string `env:"ACTION_TARGETS"`
}
