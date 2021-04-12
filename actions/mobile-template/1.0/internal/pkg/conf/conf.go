package conf

// Conf mobile template ction param collection
type Conf struct {
	WorkDir string `env:"WORKDIR"`
	// 用户指定
	ProjectName string `env:"ACTION_PROJECT_NAME" required:"false"`
	DisplayName string `env:"ACTION_DISPLAY_NAME" required:"true"`
	BundleID    string `env:"ACTION_BUNDLE_ID" required:"true"`
	PackageName string `env:"ACTION_PACKAGE_NAME" required:"true"`
}
