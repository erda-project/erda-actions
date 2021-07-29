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
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/url"
	"path/filepath"
)

type Object interface {
	Bucket() string
	Remote() string
	Local() string
}

type Client struct {
	endpoint        string
	accessKeyID     string
	accessKeySecret string

	client *oss.Client
}

// New returns the *oss.Client for uploading object to and delete object from OSS
func New(endpoint, accessKeyID, accessKeySecret string) (*Client, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		client:          client,
	}, nil
}

// Upload uploads the element to the OSS
func (o *Client) Upload(obj Object) (string, error) {
	bucket, err := o.client.Bucket(obj.Bucket())
	if err != nil {
		return "", err
	}

	if err = bucket.PutObjectFromFile(obj.Remote(), obj.Local()); err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme: "http",
		Host:   obj.Bucket() + "." + o.endpoint,
		Path:   filepath.Join("/", obj.Remote()),
	}

	return u.String(), nil
}

// DeleteRemote deletes the element from the OSS
func (o *Client) DeleteRemote(ele Object) error {
	bucket, err := o.client.Bucket(ele.Bucket())
	if err != nil {
		return err
	}

	return bucket.DeleteObject(ele.Remote())
}
