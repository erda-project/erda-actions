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
	"github.com/erda-project/erda/apistructs"
)

type ReleasesResponse struct {
	apistructs.Header
	Data *ReleasesResponseData `json:"data" yaml:"data"`
}
type ReleasesResponseData struct {
	Total uint64                 `json:"total" yaml:"total"`
	List  []ReleasesResponseItem `json:"list" yaml:"list"`
}

type ReleasesResponseItem struct {
	ReleaseID       string `json:"releaseId" yaml:"releaseId"`
	Version         string `json:"version" yaml:"version"`
	ApplicationName string `json:"applicationName" yaml:"applicationName"`
	Labels          Labels `json:"labels" yaml:"labels"`
	IsLatest        bool   `json:"isLatest" yaml:"isLatest"`
	IsStable        bool   `json:"isStable" yaml:"isStable"`
}

type Labels struct {
	GitBranch        string `json:"gitBranch"`
	GitCommitId      string `json:"gitCommitId"`
	GitCommitMessage string `json:"gitCommitMessage"`
	GitRepo          string `json:"gitRepo"`
}

type CreateUpdateReleaseRequest struct {
	Version                string     `json:"version"`
	ApplicationReleaseList [][]string `json:"applicationReleaseList"`
	Changelog              string     `json:"changelog"`
	OrgId                  uint64     `json:"orgId"`
	UserId                 string     `json:"userId"`
	ProjectID              int64      `json:"projectID"`
	IsStable               bool       `json:"isStable"`
	IsFormal               bool       `json:"isFormal"`
	IsProjectRelease       bool       `json:"isProjectRelease"`
}
