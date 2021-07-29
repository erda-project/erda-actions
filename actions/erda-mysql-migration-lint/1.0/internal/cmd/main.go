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
	"github.com/erda-project/erda-actions/pkg/metawriter"
	"github.com/erda-project/erda/pkg/database/sqlparser/migrator"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/erda-project/erda-actions/actions/erda-mysql-migration-lint/1.0/internal/config"
)

func main() {
	logrus.Infoln("Erda MySQL Migration Lint start working")
	logrus.Infof("Configuration:\n%s", config.ConfigurationString())

	var (
		err error
		c   = config.Configuration()
	)

	defer func() {
		_ = metawriter.Write(map[string]interface{}{"success": err == nil, "error": err})
	}()

	scripts, err := migrator.NewScripts(c)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to load scripts")
	}

	scripts.IgnoreMarkPending()

	if err = scripts.SameNameLint(); err != nil {
		logrus.Fatalln(err)
	}

	if err = scripts.AlterPermissionLint(); err != nil {
		logrus.Fatalln(err)
	}

	if err = scripts.Lint(); err != nil {
		logrus.Fatalln(err)
	}

	logrus.Println("Erda MySQL Migration Lint OK")

	os.Exit(0)
}
