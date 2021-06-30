package build

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/erda-project/erda/pkg/http/customhttp"
)

type HttpResponse struct {
	Success bool   `json:"success,omitempty"`
	Err     ErrMsg `json:"err,omitempty"`
}

type ErrMsg struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type ApiCheck struct {
	OrgId       string      `json:"orgId,omitempty"`
	ProjectId   string      `json:"projectId,omitempty"`
	Workspace   string      `json:"workspace,omitempty"`
	ClusterName string      `json:"clusterName,omitempty"`
	AppId       string      `json:"appId,omitempty"`
	RuntimeName string      `json:"runtimeName,omitempty"`
	ServiceName string      `json:"serviceName,omitempty"`
	Swagger     interface{} `json:"swagger,omitempty"`
}

func DoRequest(client *http.Client, method, url string, body []byte, timeout int, headers ...map[string]string) ([]byte, *http.Response, error) {
	client.Timeout = time.Duration(timeout) * time.Second
	respBody := []byte("")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	req, err := customhttp.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return respBody, nil, err
	}
	for _, kv := range headers {
		for key, value := range kv {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return respBody, nil, err
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return respBody, resp, nil
}

func Request(method, url string, body []byte, timeout int, headers ...map[string]string) ([]byte, *http.Response, error) {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return DoRequest(client, method, url, body, timeout, headers...)
}
