package oapi

import (
	"fmt"
	"net/http"
)

type AccessAPI struct {
	host      string
	token     string
	org       string
	releaseID string
}

func NewAccessAPI(host, token, org, releaseID string) *AccessAPI {
	return &AccessAPI{
		host:      host,
		token:     token,
		org:       org,
		releaseID: releaseID,
	}
}

func (api *AccessAPI) RequestHeader() http.Header {
	return http.Header{
		"authorization": {api.token},
		"org-id":        {api.org},
		"content-type":  {"application/json"},
	}
}

func (api *AccessAPI) GetDiceURL() string {
	return api.host + fmt.Sprintf("/api/releases/%s/actions/get-dice", api.releaseID)
}
