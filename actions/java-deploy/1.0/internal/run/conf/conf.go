package conf

import (
	"github.com/erda-project/erda/pkg/envconf"
)

var cfg conf

type conf struct {
	userConf     userConf
	platformConf platformConf
}

type userConf struct {
	Workdir   string `env:"ACTION_WORKDIR" required:"true"`
	Registry  string `env:"ACTION_REGISTRY" required:"true"`
	Username  string `env:"ACTION_USERNAME" required:"true"`
	Password  string `env:"ACTION_PASSWORD" required:"true"`
	SkipTests bool   `env:"ACTION_SKIP_TESTS" default:"true"` // 跳过测试
	Modules   string `env:"ACTION_MODULES"`                   // 逗号分隔 => -am -pl ${MODULES}
	Cmd       string `env:"ACTION_CMD"`                       // 用户指定的发布命令
}

type platformConf struct {
	MetaFile             string `env:"METAFILE" required:"true"`
	WorkDir              string `env:"WORKDIR" required:"true"`
	ClusterNexusURL      string `env:"BP_NEXUS_URL" required:"true"`      // 集群 nexus url，org 级别 nexus group 上线后会注入该地址
	ClusterNexusUsername string `env:"BP_NEXUS_USERNAME" required:"true"` // nexus url 对应的 username
	ClusterNexusPassword string `env:"BP_NEXUS_PASSWORD" required:"true"` // nexus url 对应的 password
}

func LoadEnvConfig() error {
	if err := envconf.Load(&cfg.userConf); err != nil {
		return err
	}
	if err := envconf.Load(&cfg.platformConf); err != nil {
		return err
	}
	return nil
}

func UserConf() userConf {
	return cfg.userConf
}

func PlatformConf() platformConf {
	return cfg.platformConf
}
