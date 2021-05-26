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

package oss

import (
	"github.com/erda-project/erda/pkg/cloudstorage"
)

type OSSEle interface {
	Bucket() string
	Remote() string
	Local() string
}

func New(endpoint, key, secret string) (*Uploader, error) {
	client, err := cloudstorage.New(endpoint, key, secret)
	if err != nil {
		return nil, err
	}
	return &Uploader{client: client}, nil
}

type Uploader struct {
	client cloudstorage.Client
}

func (u *Uploader) Upload(ele OSSEle) error {
	_, err := u.client.UploadFile(ele.Bucket(), ele.Remote(), ele.Local())
	return err
}
