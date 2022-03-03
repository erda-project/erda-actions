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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/project-artifacts/1.0/internal/config"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

var list *releases

type releases struct {
	cfg  *config.Config
	list map[string]ReleasesResponseItem
}

func singleton(cfg *config.Config) error {
	if list == nil {
		list = &releases{
			cfg:  cfg,
			list: make(map[string]ReleasesResponseItem),
		}
		if err := getApplicationsReleases(cfg); err != nil {
			return errors.Wrap(err, "failed to getApplicationsReleases")
		}
	}
	return nil
}

func (r *releases) getByReleaseID(cfg *config.Config, releaseID string) (string, bool, error) {
	item, ok := r.list[releaseID]
	if ok {
		return item.ReleaseID, true, nil
	}
	return getApplicationReleaseByID(r.cfg, releaseID)
}

func (r *releases) getByBranch(appName, branch string) (string, bool) {
	for releaseID, release := range r.list {
		if release.IsStable && release.IsLatest && release.ApplicationName == appName && release.Labels.GitBranch == branch {
			return releaseID, true
		}
	}
	return "", false
}

func getApplicationsReleases(cfg *config.Config) error {
	// var request and response
	var resp ReleasesResponse
	response, err := httpclient.New().
		Get(cfg.OapiHost).
		Path(oAPI).
		Param("pageSize", "65535").
		Param("projectId", strconv.FormatInt(cfg.ProjectID, 10)).
		Param("isProjectRelease", strconv.FormatBool(false)).
		Header("Authorization", cfg.OapiToken).
		Do().
		JSON(&resp)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}
	if !resp.Success {
		return errors.Errorf("failed to do request: %+v, data: %s", resp.Error, string(response.Body()))
	}
	if resp.Data == nil || len(resp.Data.List) == 0 {
		return nil
	}
	debugf("total: %v, this page: %v", resp.Data.Total, len(resp.Data.List))
	for i, item := range resp.Data.List {
		debugf("name: %s, branch: %s, releaseID: %s, isStable: %v, isLatest: %v",
			item.ApplicationName, item.Labels.GitBranch, item.ReleaseID, item.IsStable, item.IsLatest)
		list.list[item.ReleaseID] = resp.Data.List[i]
	}
	return nil
}

func getApplicationReleaseByID(cfg *config.Config, releaseID string) (string, bool, error) {
	// var request and response
	var resp apistructs.Header
	if _, err := httpclient.New().
		Get(cfg.OapiHost).Path(filepath.Join(oAPI, releaseID)).
		Header("Authorization", cfg.OapiToken).
		Do().
		JSON(&resp); err != nil {
		return "", false, errors.Wrap(err, "failed to request")
	}
	if !resp.Success {
		return "", false, errors.Errorf("releaseID %s not found", releaseID)
	}
	return releaseID, true, nil
}

func debugf(format string, args ...interface{}) {
	if strings.EqualFold(os.Getenv("PIPELINE_DEBUG"), "true") {
		logrus.Debugf(format, args...)
	}
}
