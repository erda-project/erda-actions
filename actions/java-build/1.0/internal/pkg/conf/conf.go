package conf

type Conf struct {
	MetaFile    string `env:"METAFILE"`
	WorkDir     string `env:"WORKDIR" default:"."`
	DiceVersion string `env:"DICE_VERSION"`
	// 用户指定
	ServiceName    string      `env:"ACTION_SERVICE_NAME"`
	BuildCmd       []string    `env:"ACTION_BUILD_CMD" required:"true"`
	Context        string      `env:"ACTION_WORKDIR" required:"true"` // maven package 目录
	JDKVersion     interface{} `env:"ACTION_JDK_VERSION" `
	Service        string      `env:"ACTION_SERVICE"`                   // Deprecated: 与 dice.yml 里 service 对应，部署时，通过 service 关联镜像
	Profile        string      `env:"ACTION_PROFILE" default:"default"` // Deprecated: spring.profiles.active
	MonitorAgent   string      `env:"ACTION_MONITOR" default:"true"`    // 是否使用监控 agent，若用户未配置，默认启用, true/false
	PreStartScript string      `env:"ACTION_PRE_START_SCRIPT"`          // 执行用户运行前脚本路径+名称，默认为项目根目录
	PreStartArgs   string      `env:"ACTION_PRE_START_ARGS"`            // 执行用户运行前脚本参数
	// pipeline注入，镜像生成需要
	OrgID             int64  `env:"DICE_ORG_ID" required:"true"`
	OrgName           string `env:"DICE_ORG_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	AppID             int64  `env:"DICE_APPLICATION_ID" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE"`
	TaskName          string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ClusterName       string `env:"DICE_CLUSTER_NAME" required:"true"`
	GittarRepo        string `env:"GITTAR_REPO"`
	GittarBranch      string `env:"GITTAR_BRANCH"`
	ProjectAppAbbr    string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceOperatorId    string `env:"DICE_OPERATOR_ID" default:"terminus"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	// pipeline注入，集群级别配置
	CentralRegistry       string `env:"BP_DOCKER_BASE_REGISTRY"`     // 中心集群 registry, eg: registry.erda.cloud
	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`
	NexusAddr             string `env:"BP_NEXUS_URL"`
	NexusUsername         string `env:"BP_NEXUS_USERNAME"`
	NexusPassword         string `env:"BP_NEXUS_PASSWORD"`
	// pipeline注入，docker build资源限制
	CPU    float64 `env:"PIPELINE_LIMITED_CPU" default:"0.5"`  // 核数, eg: 0.5
	Memory int     `env:"PIPELINE_LIMITED_MEM" default:"2048"` // 单位: M
}
