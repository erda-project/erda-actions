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

// read migrations scripts
package repo

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/config"
)

type Script struct {
	NameFromService string
	Filename        string

	conf *config.Config
}

func (s *Script) Bucket() string {
	return s.conf.OssBucket
}

// Remote is like /archived-versions/{git-tag:v1.0.0}/sqls/{service-name:cmdb}/{filename:20210101-01-base.sql}
func (s *Script) Remote() string {
	return filepath.Join(s.conf.GetOssPath(), "sqls", s.NameFromService)
}

func (s *Script) Local() string {
	return s.Filename
}

// ReadScripts read all migrations scripts.
// if the workdir or migDir is empty, it read nothing.
func ReadScripts(conf *config.Config) ([]*Script, error) {
	if conf.Workdir == "" || conf.MigDir == "" {
		return nil, nil
	}

	migdir := filepath.Join(conf.Workdir, conf.MigDir)
	services, err := ioutil.ReadDir(migdir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ReadDir %s", migdir)
	}

	var scripts []*Script
	for _, service := range services {
		logrus.Debugln("service name:", service.Name())

		if !service.IsDir() {
			continue
		}

		serviceDir := filepath.Join(conf.Workdir, conf.MigDir, service.Name())
		files, err := ioutil.ReadDir(serviceDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to ReadDir %s", serviceDir)
		}

		for _, file := range files {
			logrus.Debugln("\tfile name:", file.Name())

			if file.IsDir() {
				continue
			}

			var script = Script{
				NameFromService: filepath.Join(service.Name(), file.Name()),
				Filename:        filepath.Join(conf.Workdir, conf.MigDir, service.Name(), file.Name()),
				conf:            conf,
			}

			scripts = append(scripts, &script)
		}
	}

	return scripts, nil
}

func (s Script) Module() Module {
	return Module{
		bucket: s.Bucket(),
		remote: filepath.Dir(s.Remote()),
	}
}

type Module struct {
	bucket string
	remote string
}

func (m Module) Bucket() string {
	return m.bucket
}

func (m Module) Remote() string {
	return m.remote
}
