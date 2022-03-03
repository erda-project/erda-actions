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

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/project-artifacts/1.0/internal/config"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

var (
	oAPI = "/api/releases"
)

// GetReleaseID gets releaseID for given app conditions
func GetReleaseID(cfg *config.Config, app config.Application) (string, bool, error) {
	if err := singleton(cfg); err != nil {
		return "", false, err
	}
	if app.ReleaseID != "" {
		return list.getByReleaseID(cfg, app.ReleaseID)
	}
	releaseID, ok := list.getByBranch(app.Name, app.Branch)
	return releaseID, ok, nil
}

// CreateProjectRelease creates the project release if it is not created yet,
// updates the project release if it is already created.
func CreateProjectRelease(cfg *config.Config, releases [][]string) (string, error) {
	var request = CreateUpdateReleaseRequest{
		Version:                cfg.Version,
		ApplicationReleaseList: releases,
		Changelog:              cfg.ChangeLog,
		OrgId:                  cfg.OrgID,
		UserId:                 cfg.UserID,
		ProjectID:              cfg.ProjectID,
		IsStable:               true,
		IsFormal:               false,
		IsProjectRelease:       true,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", errors.Wrapf(err, "failed to Marshal: %+v", request)
	}
	var resp = struct {
		apistructs.Header
		Data ReleasesResponseItem
	}{}
	response, err := httpclient.New().Post(cfg.OapiHost).
		Path(oAPI).
		Header("Authorization", cfg.OapiToken).
		RawBody(bytes.NewReader(body)).
		Do().
		JSON(&resp)
	if err != nil {
		return "", errors.Wrap(err, "failed to request")
	}
	if !resp.Success {
		return "", errors.Errorf("failed to request: %+v, data: %s", resp.Error, string(response.Body()))
	}
	if resp.Data.ReleaseID == "" {
		return "", errors.New("bad response: releaseID is empty")
	}
	return resp.Data.ReleaseID, nil
}
