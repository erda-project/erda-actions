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

package client

import (
	"time"
)

type StatusInfo struct {
	UserInfo

	SessionID string     `json:"sessionid"`
	ExpiredAt *time.Time `json:"expiredAt"`
	OrgID     uint64     `json:"orgID"`
}

type UserInfo struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	NickName    string `json:"nickName"`
	Enabled     bool   `json:"enabled"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	LastLoginAt string `json:"lastLoginAt"`
}
