package conf

// conf js action param collection
type Conf struct {
	MetaFile string `env:"METAFILE"`
	WorkDir  string `env:"WORKDIR"`
	// user action params
	OutParams []OutParam `env:"ACTION_OUT_PARAMS"`
	Data      string     `env:"ACTION_DATA"`
}

type OutParam struct {
	Key        string `json:"key"`
	Expression string `json:"expression"`
}
