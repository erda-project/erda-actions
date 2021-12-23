package conf

type Conf struct {
	DiceVersion string `env:"DICE_VERSION"`
	// 用户指定
	EmailTemplateAddr   string   `env:"ACTION_EMAIL_TEMPLATE_ADDR"`
	EmailTemplateObject string   `env:"ACTION_EMAIL_TEMPLATE_OBJECT"`
	ToMail              []string `env:"ACTION_TO_EMAIL"`
	// pipeline注入，镜像生成需要
	OrgID       int64  `env:"DICE_ORG_ID" required:"true"`
	OrgName     string `env:"DICE_ORG_NAME" required:"true"`
	ProjectID   int64  `env:"DICE_PROJECT_ID" required:"true"`
	ProjectName string `env:"DICE_PROJECT_NAME" required:"true"`
	AppID       int64  `env:"DICE_APPLICATION_ID" required:"true"`
	AppName     string `env:"DICE_APPLICATION_NAME" required:"true"`
	Workspace   string `env:"DICE_WORKSPACE"`
	TaskName    string `env:"PIPELINE_TASK_NAME" default:"unknown"`
	ClusterName string `env:"DICE_CLUSTER_NAME" required:"true"`
}
