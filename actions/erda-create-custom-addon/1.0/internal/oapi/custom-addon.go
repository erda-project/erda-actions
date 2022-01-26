// Copyright (c) 2022 Terminus, Inc.
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

package oapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/erda-create-custom-addon/1.0/internal/config"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type CustomAddon struct {
	cfg *config.Config
	req *CreateCustomAddonBody
}

func New(c *config.Config) *CustomAddon {
	return &CustomAddon{
		cfg: c,
		req: &CreateCustomAddonBody{
			AddonName:       "custom",
			CreateType:      "create",
			Name:            c.Name,
			Workspace:       strings.ToUpper(c.Workspace),
			Tag:             c.Tag,
			CustomAddonType: "custom",
			Configs:         c.GetConfigs(),
			ProjectId:       c.ProjectID,
		},
	}
}

func (a *CustomAddon) Create() error {
	var l = logrus.WithField("func", "*CustomAddon.Create")

	if _, err := a.Get(); err == nil {
		l.WithField("name", a.req.Name).
			WithField("workspace", a.req.Workspace).
			Infoln("the addon exists, skip creating")
		return nil
	}

	// to create the custom addon
	var body = bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(a.req); err != nil {
		l.WithError(err).Errorln("failed to Encode")
		return errors.Wrap(err, "failed to Encode")
	}
	var resp struct {
		header
		Data *CreateCustomAddonResponseData `json:"data"`
	}
	if _, err := httpclient.New().Post(a.cfg.OpenapiHost).
		Path("/api/addons/actions/create-custom").
		Header("Authorization", a.cfg.OpenapiToken).
		RawBody(body).
		Do().
		JSON(&resp); err != nil {
		l.WithError(err).Errorln("failed to Post or Do JSON")
		return err
	}
	if !resp.Success {
		l.Errorf("failed to create addon: %+v", resp.Err)
		return errors.Errorf("failed to create addon: %+v", resp.Err)
	}
	return nil
}

func (a *CustomAddon) Get() (*AddonInfo, error) {
	var l = logrus.WithField("func", "*CustomAddon.Get")

	// list all of addons, if the same name custom addon exists, return it
	addons, err := a.list()
	if err != nil {
		l.WithError(err).Errorln("failed to list addons")
		return nil, errors.Wrap(err, "failed to list addons")
	}
	if addon, ok := sameNameIn(addons, a.req.Name, a.req.Workspace); ok {
		return a.get(addon.InstanceId)
	}
	return nil, errors.Errorf("custom addon not found, name: %s, workspace: %s", a.req.Name, a.req.Workspace)
}

func (a *CustomAddon) list() ([]*AddonInfo, error) {
	var l = logrus.WithField("func", "*CustomAddon.list")
	var resp struct {
		header
		Data []*AddonInfo `json:"data"`
	}
	if _, err := httpclient.New().Get(a.cfg.OpenapiHost).
		Path("/api/addons").
		Params(url.Values{
			"type":  {"project"},
			"value": {strconv.FormatInt(a.cfg.ProjectID, 10)},
		}).
		Header("Authorization", a.cfg.OpenapiToken).
		Do().
		JSON(&resp); err != nil {
		l.WithError(err).Errorln("failed to Get or Do JSON")
		return nil, err
	}
	if !resp.Success {
		l.Errorf("failed to list addons: %+v", resp.Err)
		return nil, errors.Errorf("failed to list addons: %+v", resp.Err)
	}
	return resp.Data, nil
}

func (a *CustomAddon) get(addonInstanceID string) (*AddonInfo, error) {
	var l = logrus.WithField("func", "*CustomAddon.get")
	var resp struct {
		header
		Data *AddonInfo `json:"data"`
	}
	if _, err := httpclient.New().Get(a.cfg.OpenapiHost).
		Path(fmt.Sprintf("/api/addons/%s", addonInstanceID)).
		Header("Authorization", a.cfg.OpenapiToken).
		Do().
		JSON(&resp); err != nil {
		l.WithError(err).WithField("instanceId", addonInstanceID).Errorln("failed to Get or Do JSON")
	}
	if !resp.Success {
		l.WithField("instanceId", addonInstanceID).Errorf("failed to get addon: %+v", resp.Err)
		return nil, errors.Errorf("failed to get addon: %+v", resp.Err)
	}
	return resp.Data, nil
}

type AddonInfo struct {
	InstanceId     string          `json:"instanceId"`
	Name           string          `json:"name"`
	Tag            string          `json:"tag"`
	AddonName      string          `json:"addonName"`
	DisplayName    string          `json:"displayName"`
	Config         json.RawMessage `json:"config"`
	Cluster        string          `json:"cluster"`
	OrgId          int             `json:"orgId"`
	ProjectId      int             `json:"projectId"`
	ProjectName    string          `json:"projectName"`
	Workspace      string          `json:"workspace"`
	Status         string          `json:"status"`
	RealInstanceId string          `json:"realInstanceId"`
}

type CreateCustomAddonBody struct {
	AddonName       string            `json:"addonName"`
	CreateType      string            `json:"createType"`
	Name            string            `json:"name"`
	Workspace       string            `json:"workspace"`
	Tag             string            `json:"tag"`
	CustomAddonType string            `json:"customAddonType"`
	Configs         map[string]string `json:"configs"`
	ProjectId       int64             `json:"projectId"`
}

type CreateCustomAddonResponseData struct {
	InstanceId        string `json:"instanceId"`
	RoutingInstanceId string `json:"routingInstanceId"`
}

func sameNameIn(addons []*AddonInfo, name, workspace string) (*AddonInfo, bool) {
	for _, addon := range addons {
		if strings.EqualFold(addon.AddonName, "custom") &&
			addon.Name == name &&
			strings.EqualFold(addon.Workspace, workspace) {
			return addon, true
		}
	}
	return nil, false
}

type header struct {
	Success bool `json:"success"`
	Err     struct {
		Code string      `json:"code"`
		Msg  string      `json:"msg"`
		Ctx  interface{} `json:"ctx"`
	} `json:"err"`
}
