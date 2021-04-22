package archive

import (
	"fmt"
	"net/http"
)

type AccessAPI struct {
	host            string
	token           string
	org             string
	projectName     string
	applicationName string
	releaseID       string
}

func NewAccessAPI(host, token, org, projectName, applicationName, releaseID string) *AccessAPI {
	return &AccessAPI{
		host:            host,
		token:           token,
		org:             org,
		projectName:     projectName,
		applicationName: applicationName,
		releaseID:       releaseID,
	}
}

func (api *AccessAPI) CreateBranchURL() string {
	return api.host + fmt.Sprintf("/api/repo/%s/%s/branches", api.projectName, api.applicationName)
}

func (api *AccessAPI) CreateCommitURL() string {
	return api.host + fmt.Sprintf("/api/repo/%s/%s/commits", api.projectName, api.applicationName)
}

func (api *AccessAPI) CreateMergeRequestURL() string {
	return api.host + fmt.Sprintf("/api/repo/%s/%s/merge-requests", api.projectName, api.applicationName)
}

func (api *AccessAPI) RequestHeader() http.Header {
	return map[string][]string{
		"authorization": {api.token},
		"org-id":        {api.org},
		"content-type":  {"application/json"},
	}
}

func (api *AccessAPI) GetDiceURL() string {
	return api.host + fmt.Sprintf("/api/releases/%s/actions/get-dice", api.releaseID)
}
