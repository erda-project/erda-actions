package conf

// Conf js action param collection
type Conf struct {
	MetaFile string `env:"METAFILE"`
	WorkDir  string `env:"WORKDIR"`
	// 用户指定
	Context       string `env:"ACTION_WORKDIR" required:"true"` // npm run xxx 所在目录
	DependencyCmd string `env:"ACTION_DEPENDENCY_CMD" default:"npm ci"`
	BuildCmd      string `env:"ACTION_BUILD_CMD" required:"true"`
	ContainerType string `env:"ACTION_CONTAINER_TYPE" required:"true"`
	DestDir       string `env:"ACTION_DEST_DIR" required:"true"`
	NpmRegistry   string `env:"ACTION_NPM_REGISTRY"`
	NpmUsername   string `env:"ACTION_NPM_USER"`
	NpmPassword   string `env:"ACTION_NPM_PASSWORD"`
	Service       string `env:"ACTION_SERVICE"` // 与 dice.yml 里 service 对应，部署时，通过 service 关联镜像 TODO deprecated
	// pipeline 注入，镜像生成时使用
	TaskName       string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ClusterName    string `env:"DICE_CLUSTER_NAME" required:"true"`
	GittarRepo     string `env:"GITTAR_REPO"`
	GittarBranch   string `env:"GITTAR_BRANCH"`
	ProjectAppAbbr string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceOperatorId string `env:"DICE_OPERATOR_ID" default:"terminus"`
	DiceVersion    string `env:"DICE_VERSION"`
	// pipeline注入，集群级别配置
	CentralRegistry       string `env:"BP_DOCKER_BASE_REGISTRY"`     // 中心集群 registry, eg: registry.erda.cloud
	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`
	// pipeline注入，docker build资源限制
	CPU    float64 `env:"PIPELINE_LIMITED_CPU" default:"1"`    // 核数, eg: 1
	Memory int     `env:"PIPELINE_LIMITED_MEM" default:"2048"` // 单位: M
	BuildkitEnable string `env:"BUILDKIT_ENABLE"`
}
