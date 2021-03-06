package run

type Request struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Params struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Path     string `json:"path"`
	Tag      string `json:"tag"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
