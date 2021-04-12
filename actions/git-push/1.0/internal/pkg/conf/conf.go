package conf

// Conf git push action param collection
type Conf struct {
	WorkDir string `env:"WORKDIR"`
	// 用户指定
	Context   string `env:"ACTION_WORKDIR" required:"true"`
	RemoteUrl string `env:"ACTION_REMOTE_URL" required:"true"`
}
