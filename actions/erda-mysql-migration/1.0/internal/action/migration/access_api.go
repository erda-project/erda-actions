// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package migration

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/erda-project/erda/pkg/http/customhttp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Err     interface{}     `json:"err"`
}

type GetAddonsListResponseDataEle struct {
	InstanceID string `json:"instanceId"`
	Name       string `json:"name"`
	AddonName  string `json:"addonName"`
	OrgID      uint64 `json:"orgId"`
	ProjectID  uint64 `json:"projectId"`
	Workspace  string `json:"workspace"`
}

type GetAddonDetailResponseData struct {
	InstanceID string                           `json:"instanceId"`
	Name       string                           `json:"name"`
	AddonName  string                           `json:"addonName"`
	Workspace  string                           `json:"workspace"`
	Config     GetAddonDetailResponseDataConfig `json:"config"`
}

type GetAddonDetailResponseDataConfig struct {
	MySQLHost     string `json:"MYSQL_HOST"`
	MySQLPassword string `json:"MYSQL_PASSWORD"`
	MySQLPort     string `json:"MYSQL_PORT"`
	MySQLUserName string `json:"MYSQL_USERNAME"`
}

type GetAddonReferencesResponseData []struct {
	ApplicationID   uint64 `json:"applicationId"`
	ApplicationName string `json:"applicationName"`
	OrgID           uint64 `json:"orgId"`
	ProjectID       uint64 `json:"projectId"`
	ProjectName     string `json:"projectName"`
	RuntimeID       uint64 `json:"runtimeId"`
	RuntimeName     string `json:"runtimeName"`
}

func RequestGet(url string, timeout int, header http.Header) ([]byte, *http.Response, error) {
	client := http.Client{Timeout: time.Second * time.Duration(timeout)}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	request, err := customhttp.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	for k, values := range header {
		for _, v := range values {
			request.Header.Add(k, v)
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return body, response, nil
}

func getAddonList(url string, header http.Header) ([]GetAddonsListResponseDataEle, error) {
	body, _, err := RequestGet(url, 60, header)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to RequestGet %s", url)
	}
	logrus.Debugf("getAddonList: %s", string(body))

	var (
		response Response
		data     []GetAddonsListResponseDataEle
	)
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func getAddonDetail(url string, header http.Header) (*GetAddonDetailResponseData, error) {
	body, _, err := RequestGet(url, 60, header)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to RequestGet %s", url)
	}
	logrus.Debugf("getAddonDetail: %s", string(body))

	var (
		response Response
		data     GetAddonDetailResponseData
	)
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func getAddonReferences(url string, header http.Header) (GetAddonReferencesResponseData, error) {
	body, _, err := RequestGet(url, 60, header)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to RequestGet %s", url)
	}
	logrus.Debugf("getAddonReferences: %s", string(body))

	var (
		response Response
		data     GetAddonReferencesResponseData
	)
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}

	return data, nil
}
