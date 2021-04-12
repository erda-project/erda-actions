package conf

type Conf struct {
	WorkDir       string   `env:"WORKDIR"`
	Commands      []string `env:"ACTION_COMMANDS"`
	Target        string   `env:"ACTION_TARGET"`
	Context       string   `env:"ACTION_CONTEXT"`
	NexusUrl      string   `env:"BP_NEXUS_URL"`
	NexusUsername string   `env:"BP_NEXUS_USERNAME"`
	NexusPassword string   `env:"BP_NEXUS_PASSWORD"`
	PipelineID    string   `env:"PIPELINE_ID"`
}
