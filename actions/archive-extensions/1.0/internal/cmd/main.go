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

package main

import (
	"path/filepath"

	"github.com/erda-project/erda/pkg/cloudstorage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/archive-extensions/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/archive-extensions/1.0/internal/tar"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func main() {
	log.Init()
	logrus.Infoln("Archive Extensions action start working")

	conf, err := config.New()
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	if len(conf.Repos) == 0 {
		err = errors.New("there is no repo to archive, did you add it ?")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}
	if conf.OssEndpoint == "" || conf.OssKey == "" || conf.OssSecret == "" {
		err = errors.New("missing OSS parameters, did you set it ?")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	// tar all repos
	var tars []string
	for _, repo := range conf.Repos {
		dst, err := tar.Tar(repo, "")
		if err != nil {
			err = errors.Wrapf(err, "failed to tar.Tar")
			_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
			logrus.Fatalln(err)
		}
		tars = append(tars, dst)
	}

	oss, err := cloudstorage.New(conf.OssEndpoint, conf.OssKey, conf.OssSecret)
	if err != nil {
		err = errors.Wrapf(err, "failed to cloudstorage.New, endpoint: %s, key: %s",
			conf.OssEndpoint, conf.OssKey)
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	if conf.OssBucket == "" {
		err = errors.New("OSS bucket can not be empty")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	for _, localFile := range tars {
		remoteFile := filepath.Join(conf.GetOssPath(), filepath.Base(localFile))
		if _, err = oss.UploadFile(conf.OssBucket, remoteFile, localFile); err != nil {
			err = errors.Wrapf(err, "failed to oss.UploadFile, dst file: %s, local file: %s",
				remoteFile, localFile)
			_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
			logrus.Fatalln(err)
		}
	}

	_ = metawriter.Write(map[string]interface{}{config.Success: true})
	logrus.Infoln("Archive Extensions action complete")
}
