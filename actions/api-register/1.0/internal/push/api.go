package push

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/api-register/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/http/customhttp"
)

func GetPublishStatus(cfg conf.Conf, registerId string) (bool, error) {
	url := cfg.DiceOpenapiPrefix + "/api/gateway/registrations/" + registerId + "/status"
	headers := make(map[string]string)
	headers["Authorization"] = cfg.CiOpenapiToken
	body, _, err := Request("GET", url, nil, 5, headers)
	if err != nil {
		return false, errors.WithStack(err)
	}
	register := conf.RegisterResponse{}
	err = json.Unmarshal(body, &register)
	if err != nil {
		return false, errors.Wrapf(err, "body:%s", body)
	}
	if !register.Success {
		return false, errors.Errorf("request failed, body:%s", body)
	}
	return register.Data.Completed, nil
}

func DoRequest(client *http.Client, method, url string, body []byte, timeout int, headers ...map[string]string) ([]byte, *http.Response, error) {
	client.Timeout = time.Duration(timeout) * time.Second
	respBody := []byte("")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	req, err := customhttp.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return respBody, nil, errors.WithStack(err)
	}
	for _, kv := range headers {
		for key, value := range kv {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return respBody, nil, errors.WithStack(err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	return respBody, resp, nil
}

func Request(method, url string, body []byte, timeout int, headers ...map[string]string) ([]byte, *http.Response, error) {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return DoRequest(client, method, url, body, timeout, headers...)
}
