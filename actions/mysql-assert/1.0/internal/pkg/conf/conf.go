package conf

// Conf js action param collection
type Conf struct {
	MetaFile string `env:"METAFILE"`
	WorkDir  string `env:"WORKDIR"`
	// 用户指定
	DataSource string     `env:"ACTION_DATASOURCE"` // 数据源
	Database   string     `env:"ACTION_DATABASE"`   // 数据库名称
	Sql        string     `env:"ACTION_SQL"`        // mysql 的 命令
	OutParams  []OutParam `env:"ACTION_OUT_PARAMS"` // 出参配置
	// pipeline 注入，镜像生成时使用
	OrgID            string `env:"DICE_ORG_ID"`
	TaskName         string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ClusterName      string `env:"DICE_CLUSTER_NAME" required:"true"`
	GittarRepo       string `env:"GITTAR_REPO"`
	GittarBranch     string `env:"GITTAR_BRANCH"`
	ProjectAppAbbr   string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceOperatorId   string `env:"DICE_OPERATOR_ID" default:"terminus"`
	DiceVersion      string `env:"DICE_VERSION"`
	DiceOpenapiAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	DiceOpenapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	// pipeline注入，集群级别配置
	CentralRegistry       string `env:"BP_DOCKER_BASE_REGISTRY"`     // 中心集群 registry, eg: registry.erda.cloud
	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`
	// pipeline注入，docker build资源限制
	CPU    float64 `env:"PIPELINE_LIMITED_CPU" default:"1"`    // 核数, eg: 1
	Memory int     `env:"PIPELINE_LIMITED_MEM" default:"2048"` // 单位: M
}

type OutParam struct {
	Key        string      `json:"key"`
	Expression string      `json:"expression"`
	Assert     string      `json:"assert"`
	Value      interface{} `json:"value"`
}
