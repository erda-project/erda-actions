package conf

type HttpResponse struct {
	Success bool   `json:"success,omitempty"`
	Err     ErrMsg `json:"err,omitempty"`
}

type ErrMsg struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
