package conf

type Conf struct {
	MetaFile    string   `env:"METAFILE"`
	Username    string   `env:"ACTION_USERNAME"`
	ServiceKey  string   `env:"ACTION_SERVICE_KEY"`
	ApiKey      string   `env:"ACTION_API_KEY"`
	OrgID       string   `env:"ACTION_ORG_ID"`
	AppID       string   `env:"ACTION_APP_ID"`
	Severities  []string `env:"ACTION_SEVERITIES"`
	Status      string   `env:"ACTION_STATUS" default:"Reported"`
	Expand      string   `env:"ACTION_EXPAND" default:"vulnerability_instances"`
	AssertCount int      `env:"ACTION_ASSERT_COUNT" default:"0"`
}
