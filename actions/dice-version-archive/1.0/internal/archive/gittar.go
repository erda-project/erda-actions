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

// operations about gittar: create new branch, commit and merge request
package archive

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

const (
	ActionAdd    = "add"
	PathTypeBlob = "blob"
)

type Gittar struct {
	uri *AccessAPI
}

func NewGittar(api *AccessAPI) *Gittar {
	return &Gittar{uri: api}
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

	return d.MergeID(), nil
}

type CreateBranchPayload struct {
	// new branch name
	Name string `json:"name"`
	// src branch name
	Ref string `json:"ref"`
}

type CreateCommitPayload struct {
	// commit message
	Message string `json:"message"`
	// branch name
	Branch string `json:"branch"`
	// changes
	Actions []*CreateCommitPayloadAction `json:"actions"`
}

type CreateCommitPayloadAction struct {
	// Action is always "add"
	Action string `json:"action"`
	// file's content
	Content string `json:"content"`
	// file's path
	Path string `json:"path"`
	// PathType is always "blob"
	PathType string `json:"pathType"`
}

type CreateMergeRequestPayload struct {
	// mr title
	Title string `json:"title"`
	// mr description
	Description string `json:"description"`
	// mr processor user id
	AssigneeID string `json:"assigneeId"`
	// the branch merging from
	SourceBranch string `json:"sourceBranch"`
	// the branch merging to
	TargetBranch string `json:"targetBranch"`
	// remove source branch after merged
	RemoveSourceBranch bool `json:"removeSourceBranch"`
}

type Response struct {
	Success bool            `json:"success"`
	Err     json.RawMessage `json:"err"`
	Data    json.RawMessage `json:"data"`
}

type CreateMergeRequestResponseData struct {
	Id       uint64 `json:"id"`
	MergeID_ uint64 `json:"mergeId"`
}

func (d CreateMergeRequestResponseData) ID() string {
	return strconv.FormatUint(d.Id, 10)
}

func (d CreateMergeRequestResponseData) MergeID() string {
	return strconv.FormatUint(d.MergeID_, 10)
}
