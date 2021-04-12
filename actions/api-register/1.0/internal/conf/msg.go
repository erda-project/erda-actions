package conf

type RegisterApiMsg struct {
	OrgId       string      `json:"orgId,omitempty"`
	ProjectId   string      `json:"projectId,omitempty"`
	Workspace   string      `json:"workspace,omitempty"`
	ClusterName string      `json:"clusterName,omitempty"`
	AppId       string      `json:"appId,omitempty"`
	AppName     string      `json:"appName,omitempty"`
	RuntimeId   string      `json:"runtimeId,omitempty"`
	RuntimeName string      `json:"runtimeName,omitempty"`
	ServiceName string      `json:"serviceName,omitempty"`
	ServiceAddr string      `json:"serviceAddr,omitempty"`
	Swagger     interface{} `json:"swagger,omitempty"`
}

type HttpResponse struct {
	Success bool    `json:"success,omitempty"`
	Err     ErrMsg  `json:"err,omitempty"`
	Data    DataMsg `json:"data,omitempty"`
}

type RegisterResponse struct {
	Success bool        `json:"success,omitempty"`
	Err     ErrMsg      `json:"err,omitempty"`
	Data    RegisterMsg `json:"data,omitempty"`
}

type ErrMsg struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type DataMsg struct {
	ApiRegisterId string `json:"apiRegisterId,omitempty"`
}

type RegisterMsg struct {
	Completed bool `json:"completed,omitempty"`
}
