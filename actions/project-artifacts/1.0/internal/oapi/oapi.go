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
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/project-artifacts/1.0/internal/config"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

var (
	oAPI = "/api/releases"
)

// GetLatestApplicationRelease gets the release for the given releaseID or the latest release for the
//  given branch if the releaseID is not given.
func GetLatestApplicationRelease(cfg *config.Config, app config.Application) (string, bool, error) {
	// var request and response
	var resp ReleasesResponse
	response, err := httpclient.New().
		Get(cfg.OapiHost).
		Path(oAPI).
		Param("projectId", strconv.FormatInt(cfg.ProjectID, 10)).
		Param("isProjectRelease", strconv.FormatBool(false)).
		Header("Authorization", cfg.OapiToken).
		Do().
		JSON(&resp)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to do request")
	}
	if !resp.Success {
		return "", false, errors.Errorf("failed to do request: %+v, data: %s", resp.Error, string(response.Body()))
	}
	if resp.Data == nil || len(resp.Data.List) == 0 {
		return "", false, nil
	}
	if app.ReleaseID != "" {
		for i := range resp.Data.List {
			if item := resp.Data.List[i]; item.IsStable && item.IsLatest && item.ReleaseID == app.ReleaseID {
				return app.ReleaseID, true, nil
			}
		}
		return "", false, errors.Errorf("releaseID %s not found", app.ReleaseID)
	}
	for i := range resp.Data.List {
		if item := resp.Data.List[i]; item.IsStable && item.IsLatest &&
			item.ApplicationName == app.Name && item.Labels.GitBranch == app.Branch {
			return item.ReleaseID, true, nil
		}
	}
	return "", false, nil
}

// CreateUpdateProjectRelease creates the project release if it is not created yet,
// updates the project release if it is already created.
func CreateUpdateProjectRelease(cfg *config.Config, releases [][]string) (string, error) {
	release, ok, err := GetProjectRelease(cfg)
	if err != nil {
		return "", errors.Wrapf(err, "failed to GetProjectRelease, projectID(Name): %v(%s), version: %s",
			cfg.ProjectID, cfg.ProjectName, cfg.Version)
	}
	if ok {
		return release, updateProjectRelease(cfg, release, releases)
	}
	return createProjectRelease(cfg, releases)
}

// GetProjectRelease gets the project release with the given version,
// if it is not found, return false.
func GetProjectRelease(cfg *config.Config) (string, bool, error) {
	var resp ReleasesResponse
	if _, err := httpclient.New().Get(cfg.OapiHost).
		Path(oAPI).
		Header("Authorization", cfg.OapiToken).
		Param("projectId", strconv.FormatInt(cfg.ProjectID, 10)).
		Param("isProjectRelease", strconv.FormatBool(true)).
		Do().
		JSON(&resp); err != nil {
		return "", false, errors.Wrap(err, "failed to request")
	}
	if !resp.Success {
		return "", false, errors.Errorf("failed to request: %+v", resp.Error)
	}
	if resp.Data == nil || len(resp.Data.List) == 0 {
		return "", false, nil
	}
	for i := range resp.Data.List {
		if resp.Data.List[i].Version == cfg.Version {
			return resp.Data.List[i].ReleaseID, true, nil
		}
	}
	return "", false, nil
}

// createProjectRelease creates the project release and return it's releaseID
func createProjectRelease(cfg *config.Config, releases [][]string) (string, error) {
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

// updateProjectRelease updates the project release
func updateProjectRelease(cfg *config.Config, releaseID string, releases [][]string) error {
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
		return errors.Wrapf(err, "failed to Marshal: %+v", request)
	}
	var resp = struct {
		apistructs.Header
		Data json.RawMessage `json:"data"`
	}{}

	response, err := httpclient.New().Put(cfg.OapiHost).
		Path(filepath.Join(oAPI, releaseID)).
		RawBody(bytes.NewReader(body)).
		Header("Authorization", cfg.OapiToken).
		Do().
		JSON(&resp)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	if !resp.Success {
		return errors.Errorf("failed to request: %+v, data: %s", resp.Error, string(response.Body()))
	}
	return nil
}
