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
	"net/url"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Object interface {
	RemoteObject

	Local() string
}

type RemoteObject interface {
	Bucket() string
	Remote() string
}

type defaultObject struct {
	bucket string
	remote string
	local  string
}

func (o defaultObject) Bucket() string {
	return o.bucket
}

func (o defaultObject) Remote() string {
	return o.remote
}

func (o defaultObject) Local() string {
	return o.local
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

	if err = bucket.PutObjectFromFile(obj.Remote(), obj.Local(), oss.Progress(new(ProgressListener))); err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme: "http",
		Host:   obj.Bucket() + "." + o.endpoint,
		Path:   filepath.Join("/", obj.Remote()),
	}

	return u.String(), nil
}

// DeleteRemoteRecursively deletes the element from the OSS
func (o *Client) DeleteRemoteRecursively(obj RemoteObject) error {
	if obj.Remote() == "" || obj.Remote() == "/" {
		return errors.New("invalid path")
	}

	bucket, err := o.client.Bucket(obj.Bucket())
	if err != nil {
		return err
	}

	// list all object under the obj
	marker := oss.Marker("")
	prefix := oss.Prefix(obj.Remote())
	loopLog := logrus.WithField("bucket", obj.Bucket()).WithField("path", obj.Remote())
	return deleteDirectories(bucket, marker, prefix, loopLog)
}

func deleteFiles(bucket *oss.Bucket, marker, prefix oss.Option, log *logrus.Entry) error {
	log.Infoln("ListObjects [files]")
	page, err := bucket.ListObjects(marker, prefix)
	if err != nil {
		log.Errorln("failed to ListObjects [files]")
		return err
	}
	if len(page.Objects) == 0 {
		log.Warnln("no any more files in the path")
		return nil
	}
	var files []string
	for _, file := range page.Objects {
		files = append(files, file.Key)
	}
	log.WithField("files", files).Infoln("the files is going to be deleting")
	result, err := bucket.DeleteObjects(files)
	if err != nil {
		return err
	}
	log.WithField("result", result).Infoln("the files is deleted")
	if !page.IsTruncated {
		return nil
	}

	return deleteFiles(bucket, oss.Marker(page.NextMarker), oss.Prefix(page.Prefix), log)
}

func deleteDirectories(bucket *oss.Bucket, marker, prefix oss.Option, log *logrus.Entry) error {
	// delete files before all
	if err := deleteFiles(bucket, marker, prefix, log); err != nil {
		return err
	}

	log.Infoln("ListObjects [directories]")
	page, err := bucket.ListObjects(marker, prefix, oss.Delimiter("/"))
	if err != nil {
		log.Errorln("failed to ListObjects [directories]")
		return err
	}

	var directories []string
	for _, dir := range page.Objects {
		directories = append(directories, dir.Key)
	}
	if len(directories) > 0 {
		log.WithField("directories", directories).Infoln("the directories is going to be deleting")
	}
	for _, dir := range page.Objects {
		if err := deleteDirectories(bucket, oss.Marker(""), oss.Prefix(dir.Key), log.WithField("path", dir.Key)); err != nil {
			return err
		}
	}

	// if there are any objects is not returned, delete next page
	if page.IsTruncated {
		return deleteDirectories(bucket, oss.Marker(page.NextMarker), oss.Prefix(page.Prefix), log)
	}

	// after all, delete self
	return bucket.DeleteObject(page.Prefix)
}
