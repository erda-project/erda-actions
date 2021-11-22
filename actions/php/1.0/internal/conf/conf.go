package conf

type Conf struct {
	MetaFile   string `env:"METAFILE"`
	WorkDir    string `env:"WORKDIR"`
	IndexPath  string `env:"ACTION_INDEX_PATH" default:""`
	Context    string `env:"ACTION_CONTEXT"`
	Service    string `env:"ACTION_SERVICE"`
	PHPVersion string `env:"ACTION_PHP_VERSION" default:"7.2"`

	// pipeline注入，镜像生成需要
	TaskName       string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ClusterName    string `env:"DICE_CLUSTER_NAME" required:"true"`
	GittarRepo     string `env:"GITTAR_REPO"`
	GittarBranch   string `env:"GITTAR_BRANCH"`
	ProjectAppAbbr string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceOperatorId string `env:"DICE_OPERATOR_ID" default:"terminus"`
	// pipeline注入，集群级别配置
	CentralRegistry       string `env:"BP_DOCKER_BASE_REGISTRY" default:"registry.erda.cloud"` // 中心集群 registry
	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"`                           // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`

	// BuildKit params
	BuildkitEnable string `env:"BUILDKIT_ENABLE"`
	BuildkitdAddr  string `env:"BUILDKITD_ADDR" default:"tcp://buildkitd.default.svc.cluster.local:1234"`

	// pipeline注入，docker build资源限制
	CPU    float64 `env:"PIPELINE_LIMITED_CPU" default:"0.5"`  // 核数, eg: 0.5
	Memory int     `env:"PIPELINE_LIMITED_MEM" default:"2048"` // 单位: M
}
