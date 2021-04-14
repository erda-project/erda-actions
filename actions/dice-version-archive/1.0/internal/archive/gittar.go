// gittar 相关的操作
// commit && merge request

package archive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/erda-project/erda/pkg/customhttp"
	"github.com/pkg/errors"
)

const (
	ActionAdd    = "add"
	PathTypeBlob = "blob"
)

type Gittar struct {
	uri *AccessAPI
}

func NewGittar(host, token, org, projectName, applicationName string) *Gittar {
	return &Gittar{
		uri: &AccessAPI{
			host:            host,
			token:           token,
			org:             org,
			projectName:     projectName,
			applicationName: applicationName,
		},
	}
}

func (g Gittar) URI() *AccessAPI {
	return g.uri
}

func (g Gittar) CreateCommit(payload *CreateCommitPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body, _, err := RequestPost(g.URI().CreateCommitURL(), data, g.URI().RequestHeader())
	if err != nil {
		return errors.Wrapf(err, "failed to RequestPost %s", g.URI().CreateCommitURL())
	}
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return errors.Wrapf(err, "failed to Unmarshal CreateCommit response: %s", string(body))
	}
	if !response.Success {
		return errors.Errorf("failed to CreateCommit. %s", string(response.Err))
	}
	return nil
}

func (g Gittar) CreateBranch(payload *CreateBranchPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	body, _, err := RequestPost(g.URI().CreateBranchURL(), data, g.URI().RequestHeader())
	if err != nil {
		return errors.Wrapf(err, "failed to ReqeustPost %s", g.URI().CreateBranchURL())
	}

	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return errors.Wrapf(err, "failed to Unmarshal CreateBranch response: %s", string(body))
	}
	if !response.Success {
		return errors.Errorf("failed to CreateBranch. %s", string(response.Err))
	}

	return nil
}

func (g Gittar) CreateMergeRequest(payload *CreateMergeRequestPayload) (mrID string, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	body, _, err := RequestPost(g.URI().CreateMergeRequestURL(), data, g.URI().RequestHeader())
	if err != nil {
		return "", errors.Wrapf(err, "failed to RequestPost %s", g.URI().CreateMergeRequestURL())
	}

	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return "", errors.Wrapf(err, "failed to Unmarshal CreateMergeRequest response: %s", string(body))
	}

	if !response.Success {
		return "", errors.Errorf("failed  to CreateMergeReqeust. %s", string(response.Err))
	}

	var d CreateMergeRequestResponseData
	if err = json.Unmarshal(response.Data, &d); err != nil {
		return "", errors.Wrap(err, "failed to Unmarshal CreateMergeRequestResponseData")
	}

	return d.ID(), nil
}

type AccessAPI struct {
	host            string
	token           string
	org             string
	projectName     string
	applicationName string
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

// {"name":"feature/some-release","ref":"feature/migration-cicd"}
type CreateBranchPayload struct {
	Name string `json:"name"` // 新分支名
	Ref  string `json:"ref"`  // 源分支名
}

type CreateCommitPayload struct {
	Message string                       `json:"message"` // 提交信息
	Branch  string                       `json:"branch"`  // 分支名
	Actions []*CreateCommitPayloadAction `json:"actions"` // 修改行为
}

// {"message":"Add BaseResponse","branch":"feature/some-release","actions":[{"action":"add","content":"sdfgsad","path":"BaseResponse","pathType":"blob"}]}
type CreateCommitPayloadAction struct {
	Action   string `json:"action"`   // add
	Content  string `json:"content"`  // 文本内容
	Path     string `json:"path"`     // 修改的文件的路径
	PathType string `json:"pathType"` // blob
}

// {"title":"测试提交 mr","description":"测试提交 mr","assigneeId":"2","sourceBranch":"feature/some-release","targetBranch":"feature/migration-cicd",
// "removeSourceBranch":true}
type CreateMergeRequestPayload struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	AssigneeID         string `json:"assigneeId"`
	SourceBranch       string `json:"sourceBranch"`
	TargetBranch       string `json:"targetBranch"`
	RemoveSourceBranch bool   `json:"removeSourceBranch"`
}

type Response struct {
	Success bool            `json:"success"`
	Err     json.RawMessage `json:"err"`
	Data    json.RawMessage `json:"data"`
}

type CreateMergeRequestResponseData struct {
	Id uint64 `json:"id"`
}

func (d CreateMergeRequestResponseData) ID() string {
	return strconv.FormatUint(d.Id, 10)
}

func RequestPost(url string, payload []byte, header http.Header) ([]byte, *http.Response, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	request, err := customhttp.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	for k, values := range header {
		for _, v := range values {
			request.Header.Add(k, v)
		}
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, response, nil
}
