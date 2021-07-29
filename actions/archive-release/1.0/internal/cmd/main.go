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
	"github.com/erda-project/erda/pkg/parser/diceyml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/oss"
	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/repo"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func main() {
	log.Init()
	logrus.Infoln("Archive Release action start working")

	var err error
	defer func() {
		_ = metawriter.Write(map[string]interface{}{config.Success: err == nil, config.Err: err})
	}()

	conf, err := config.New()
	if err != nil {
		logrus.WithError(err).Fatalln("failed to load config")
	}

	// make a repo object
	r, err := repo.New(conf)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to read repo")
	}

	// make a oss handle
	client, err := oss.New(conf.OssEndpoint, conf.OssKey, conf.OssSecret)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to make an OSS client")
	}

	// delete the git ref dir before all uploading
	if err = client.DeleteRemote(conf.GitRefDir()); err != nil {
		logrus.WithError(err).WithField("path", conf.GitRefDir().Remote()).
			Fatalln("failed to remove the path from OSS")
	}

	// upload released yml, migration lint config file and SQLs scripts
	url, err := client.Upload(r.ReleaseYml)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to upload release yaml file")
	}
	if r.LintConfig.Local() != "" {
		if _, err = client.Upload(r.LintConfig); err != nil {
			_ = metawriter.Write(map[string]interface{}{config.Warn: err})
			logrus.WithError(err).Warnln("failed to upload Erda MySQL Migration Lint config file")
		}
	}
	if len(r.Scripts) == 0 {
		logrus.Warnln("no migration script will be archived because there is no workdir or migrationsDir")
	}
	for _, script := range r.Scripts {
		if _, err := client.Upload(script); err != nil {
			logrus.WithError(err).WithField("filename", script.Filename).Fatalln("failed to migration script")
		}
	}

	// write oss download url and every service's image to meta
	_ = metawriter.Write(map[string]interface{}{"erda.yml": url, "gitref": conf.GitRef})
	meta := make(map[string]interface{})
	if deployable, err := r.ReleaseYml.Deployable(); err == nil {
		var obj = new(diceyml.Object)
		if err := yaml.Unmarshal([]byte(deployable), obj); err == nil {
			for serviceName, service := range obj.Services {
				if service != nil {
					meta[serviceName] = service.Image
				}
			}
		}
	}
	_ = metawriter.Write(meta)

	logrus.Infoln("Archive Release action complete")
}
