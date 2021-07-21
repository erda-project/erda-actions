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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/erda-project/erda/pkg/database/sqllint"
	"github.com/erda-project/erda/pkg/database/sqllint/configuration"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/erda-mysql-migration-lint/1.0/internal/config"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

const (
	baseScriptLabel  = "# MIGRATION_BASE"
	baseScriptLabel2 = "-- MIGRATION_BASE"
	baseScriptLabel3 = "/* MIGRATION_BASE */"
)

func main() {
	logrus.Infoln("Erda MySQL Migration Lint start working")
	logrus.Infof("Configuration:\n%s", config.ConfigurationString())

	var (
		err     error
		msg     string
		c       = config.Configuration()
		files   = new(walk).walk(filepath.Join(c.Workdir(), c.MigrationDir()), ".sql").filenames()
		rulers  = configuration.DefaultRulers()
		lintCfg *configuration.Configuration
	)

	if c.LintConfig != "" {
		lintCfg, err = configuration.FromLocal(filepath.Join(c.Workdir(), c.LintConfig))
		if err != nil {
			msg = "failed to load lint configuration"
			_ = metawriter.Write(map[string]interface{}{"success": false, "err": err, "msg": msg})
			logrus.WithError(err).
				WithField("lint config filename", c.LintConfig).
				Fatalln(msg)
		}
		rulers, err = lintCfg.Rulers()
		if err != nil {
			msg = "failed to generate lint rulers from lint configuration"
			_ = metawriter.Write(map[string]interface{}{"success": false, "err": err, "msg": msg})
			logrus.WithError(err).
				Fatalln(msg)
		}
	}

	linter := sqllint.New(rulers...)

	for _, filename := range files {
		var data []byte
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			msg = "failed to read script file"
			_ = metawriter.Write(map[string]interface{}{"success": false, "err": err, "msg": msg, "filename": filename})
			logrus.WithError(err).WithField("msg", msg).WithField("filename", filename).Fatalln(msg)
		}
		if !c.LintBase && isBaseScript(data) {
			continue
		}
		if err = linter.Input(data, filename); err != nil {
			msg = "failed to input script text to linter"
			_ = metawriter.Write(map[string]interface{}{"success": false, "err": err, "msg": msg, "filename": filename})
			logrus.WithError(err).WithField("msg", msg).WithField("filename", filename).Fatalln(msg)
		}
	}

	if len(linter.Errors()) == 0 {
		msg = "Erda MySQL Migration Lint ok"
		_ = metawriter.Write(map[string]interface{}{"success": true, "msg": msg})
		os.Exit(0)
	}

	msg = "some errors in your migrations"
	_ = metawriter.Write(map[string]interface{}{"success": false, "msg": msg})

	out := io.MultiWriter(os.Stdout, os.Stderr)
	if _, err = fmt.Fprintln(out, linter.Report()); err != nil {
		logrus.WithError(err).Fatalln("failed to print report")
	}
	for src, errs := range linter.Errors() {
		if _, err := fmt.Fprintln(out, src); err != nil {
			logrus.WithError(err).Fatalln("failed to print error")
		}
		for _, e := range errs {
			if _, err := fmt.Fprintln(out, e); err != nil {
				logrus.WithError(err).Fatalln("failed to print error")
			}
		}
	}

	logrus.Fatalln(msg)
}

type walk struct {
	files []string
}

func (w *walk) filenames() []string {
	return w.files
}

func (w *walk) walk(input, suffix string) *walk {
	infos, err := ioutil.ReadDir(input)
	if err != nil {
		w.files = append(w.files, input)
		return w
	}

	for _, info := range infos {
		if info.IsDir() {
			w.walk(filepath.Join(input, info.Name()), suffix)
			continue
		}
		if strings.ToLower(path.Ext(info.Name())) == strings.ToLower(suffix) {
			file := filepath.Join(input, info.Name())
			w.files = append(w.files, file)
		}
	}

	return w
}

func isBaseScript(data []byte) bool {
	return bytes.HasPrefix(data, []byte(baseScriptLabel)) ||
		bytes.HasPrefix(data, []byte(baseScriptLabel2)) ||
		bytes.HasPrefix(data, []byte(baseScriptLabel3))
}
