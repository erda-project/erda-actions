package build

const GoMODBuild = "gomod"
const GoVendorBuild = "govendor"
const OtherBuild = "other"

type GoBuild struct {
	BuildPath string
	BuildType string
	Package   string
	Name      string
}

type GoVendorConfig struct {
	RootPath string `json:"rootPath"`
	Comment  string `json:"comment"`
	Ignore   string `json:"ignore"`
}
