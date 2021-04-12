package pack

// ModuleImage build action构建完成后写镜像文件格式，供 release action 读取
type ModuleImage struct {
	ModuleName string `json:"module_name"`
	Image      string `json:"image"`
}
