package conf

// Conf dockerfile action param collection
type Conf struct {
	MetaFile string `env:"METAFILE"`
	WorkDir  string `env:"WORKDIR"`
	// 用户指定
	Context      string `env:"ACTION_WORKDIR" required:"true"` // dockerfile 构建目录
	Path         string `env:"ACTION_PATH" required:"true"`
	BuildArgsStr string `env:"ACTION_BUILD_ARGS"` // 用于渲染 dockerfile
	BuildArgs    map[string]string
	Service      string          `env:"ACTION_SERVICE"` // TODO deprecated
	Image        *DockerImage    `env:"ACTION_IMAGE"`
	Registry     *DockerRegistry `env:"ACTION_REGISTRY"`
	// pipeline 注入，镜像生成时使用
	TaskName       string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ProjectAppAbbr string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	DiceWorkspace  string `env:"DICE_WORKSPACE" required:"true"`
	DiceOperatorId string `env:"DICE_OPERATOR_ID" default:"terminus"`

	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`

	// BuildKit params
	BuildkitEnable string `env:"BUILDKIT_ENABLE"`
	BuildkitdAddr  string `env:"BUILDKITD_ADDR" default:"tcp://buildkitd.default.svc.cluster.local:1234"`

	// pipeline注入，docker build资源限制
	CPU    float64 `env:"PIPELINE_LIMITED_CPU" default:"1"`    // 核数, eg: 1
	Memory int     `env:"PIPELINE_LIMITED_MEM" default:"2048"` // 单位: M
}

type DockerRegistry struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DockerImage struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}
