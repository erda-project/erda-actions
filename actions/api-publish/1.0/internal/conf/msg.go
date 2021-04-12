package conf

type PublishMsg struct {
	OwnerEmail string `json:"ownerEmail,omitempty"`
	ItemName   string `json:"itemName,omitempty"`
	OrgId      string `json:"orgId,omitempty"`
}

type HttpResponse struct {
	Success bool   `json:"success,omitempty"`
	Err     ErrMsg `json:"err,omitempty"`
}

type ErrMsg struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
