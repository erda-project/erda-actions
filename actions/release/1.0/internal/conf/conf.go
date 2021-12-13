package conf

// Conf action 入参
type Conf struct {
	WorkDir  string `env:"WORKDIR"`
	Metafile string `env:"METAFILE"`
	// Params
	DiceYaml               string `env:"ACTION_DICE_YML"`
	DiceDevelopmentYaml    string `env:"ACTION_DICE_DEVELOPMENT_YML"`
	DiceTestYaml           string `env:"ACTION_DICE_TEST_YML"`
	DiceStagingYaml        string `env:"ACTION_DICE_STAGING_YML"`
	DiceProductionYaml     string `env:"ACTION_DICE_PRODUCTION_YML"`
	ReplacementImageStr    string `env:"ACTION_REPLACEMENT_IMAGES"` // <bp-name>/pack-result, pack-result存储镜像
	ReplacementImages      []string
	InitSQL                string      `env:"ACTION_INIT_SQL"` // eg: repo/db
	ReleaseFiles           string      `env:"ACTION_RELEASE_FILES"`
	ReleaseMobile          *MobileData `env:"ACTION_RELEASE_MOBILE"`
	ImageStr               string      `env:"ACTION_IMAGE"`
	Images                 map[string]string
	LabelStr               string `env:"ACTION_LABELS"`
	Labels                 map[string]string
	ServicesStr            string `env:"ACTION_SERVICES"`
	Services               map[string]Service
	CheckDiceyml           bool   `env:"ACTION_CHECK_DICEYML" default:"true"`
	TagVersion             string `env:"ACTION_TAG_VERSION"`
	MigrationType          string `env:"ACTION_MIGRATION_TYPE"`
	MigrationDir           string `env:"ACTION_MIGRATION_DIR"`
	MigrationMysqlDatabase string `env:"ACTION_MIGRATION_MYSQL_DATABASE"`
	CrossCluster           bool   `env:"ACTION_CROSS_CLUSTER" default:"false"`
	AABInfoStr             string `env:"ACTION_AAB_INFO"`
	AABInfo                AABInfo

	// env
	OrgID              int64  `env:"DICE_ORG_ID" required:"true"`
	OrgName            string `env:"DICE_ORG_NAME" required:"true"`
	ClusterName        string `env:"DICE_CLUSTER_NAME" required:"true"`
	ProjectID          int64  `env:"DICE_PROJECT_ID" required:"true"`
	ProjectName        string `env:"DICE_PROJECT_NAME" required:"true"`
	AppID              int64  `env:"DICE_APPLICATION_ID" required:"true"`
	AppName            string `env:"DICE_APPLICATION_NAME" required:"true"`
	Workspace          string `env:"DICE_WORKSPACE"`
	GittarRepo         string `env:"GITTAR_REPO"`
	GittarBranch       string `env:"GITTAR_BRANCH"`
	GittarCommitID     string `env:"GITTAR_COMMIT"`
	GittarMessage      string `env:"GITTAR_MESSAGE"`
	DiceOperatorID     string `env:"DICE_OPERATOR_ID" required:"true"`
	CiOpenapiToken     string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	PipelineStorageURL string `env:"PIPELINE_STORAGE_URL"`
	ReleaseTag         string `env:"RELEASE_TAG"`
	ProjectAppAbbr     string `env:"DICE_PROJECT_APPLICATION"` // 用于生成用户镜像repo
	TaskName           string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	DiceOperatorId     string `env:"DICE_OPERATOR_ID" default:"terminus"`
	PipelineID         string `env:"PIPELINE_ID"`

	LocalRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"` // 集群内 registry
	LocalRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	LocalRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`

	DiceVersion  string `env:"DICE_VERSION"`
	Base64Switch bool   `env:"BASE64_SWITCH"` // base64 开关
	BuildkitEnable string `env:"BUILDKIT_ENABLE"`
}

type Service struct {
	Name  string
	Cmd   string   `json:"cmd"`
	Image string   `json:"image"`
	Cps   []string `json:"copys"`
}

type MobileData struct {
	Files   []string `json:"files"`
	Version string   `json:"version"`
}

type AABInfo struct {
	PackageName interface{} `json:"packageName"`
	VersionCode interface{} `json:"versionCode"`
	VersionName interface{} `json:"versionName"`
}
