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
	"github.com/sirupsen/logrus"

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
		if err != nil {
			_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		}
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

	// make an oss handler
	client, err := oss.New(conf.OssEndpoint, conf.OssKey, conf.OssSecret)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to make an OSS client")
	}

	// upload released yml
	url, err := client.Upload(r.ReleaseYml)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to upload release yaml file")
	}

	// if there is a migration lint configuration file, upload it
	if r.LintConfig.Local() != "" {
		if _, err = client.Upload(r.LintConfig); err != nil {
			_ = metawriter.Write(map[string]interface{}{config.Warn: err})
			logrus.WithError(err).Warnln("failed to upload Erda MySQL Migration Lint config file")
		}
	}

	// if there is migration scripts, upload them
	if len(r.Scripts) == 0 {
		logrus.Warnln("no migration script will be archived because there is no workdir or migrationsDir")
	}
	var deletedModules = make(map[string]repo.Module)
	for _, script := range r.Scripts {
		// before the first time upload the script, must remove the module
		module := script.Module()
		if _, ok := deletedModules[module.Remote()]; !ok {
			if err := client.DeleteRemoteRecursively(module); err != nil {
				logrus.WithError(err).WithField("bucket", module.Bucket()).WithField("path", module.Remote()).
					Fatalln("failed to remove the path from OSS")
			}
			deletedModules[module.Remote()] = module
		}
		if _, err = client.Upload(script); err != nil {
			logrus.WithError(err).WithField("filename", script.Filename).Fatalln("failed to migration script")
		}
	}

	// write oss download url and every service's image to meta
	meta := map[string]interface{}{"erda.yml": url, "gitref": conf.GitRef, config.Success: true}
	logrus.Infoln("erda.yml", url)
	logrus.Infoln("gitref", conf.GitRef)
	if obj := r.ReleaseYml.Obj(); obj != nil {
		for serviceName, service := range obj.Services {
			logrus.Infoln(serviceName, service.Image)
			meta[serviceName] = service.Image
		}
	}
	if err := metawriter.Write(meta); err != nil {
		logrus.WithFields(meta).Warnln("failed to write info to meta")
	}

	logrus.Infoln("Archive Release action complete")
}
