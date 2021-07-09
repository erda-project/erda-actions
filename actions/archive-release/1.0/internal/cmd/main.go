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
	"github.com/pkg/errors"
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

	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "failed to config.New")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	// make a repo object
	r, err := repo.New(conf)
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	// make a oss handle
	uploader, err := oss.New(conf.OssEndpoint, conf.OssKey, conf.OssSecret)
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	// upload released yml, migration lint config file and SQLs scripts
	url, err := uploader.Upload(r.ReleaseYml)
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}
	if r.LintConfig.Local() != "" {
		if _, err = uploader.Upload(r.LintConfig); err != nil {
			_ = metawriter.Write(map[string]interface{}{config.Warn: err})
			logrus.Warnln(err)
		}
	}
	if len(r.Scripts) == 0 {
		logrus.Warnln("no migration script will be archived because there is no workdir or migrationsDir")
	}
	for _, script := range r.Scripts {
		if _, err := uploader.Upload(script); err != nil {
			_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
			logrus.Fatalln(err)
		}
	}

	// write oss download url and every service's image to meta
	meta := map[string]interface{}{config.Success: true, "erda.yml": url, "gitref": conf.GitRef}
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
