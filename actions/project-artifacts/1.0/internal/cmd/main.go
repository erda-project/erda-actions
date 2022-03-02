// Copyright (c) 2022 Terminus, Inc.
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
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/project-artifacts/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/project-artifacts/1.0/internal/oapi"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func main() {
	infoL := logrus.New()
	infoL.SetOutput(os.Stdout)
	errL := logrus.New()
	errL.SetOutput(os.Stderr)

	cfg, err := config.GetConfig()
	if err != nil {
		_ = metawriter.WriteSuccess(false)
		_ = metawriter.WriteError(err)
		errL.WithError(err).Fatalf("failed to GetConfig")
	}
	cfg.Print()
	groups, err := cfg.GetGroups()
	if err != nil {
		_ = metawriter.WriteSuccess(false)
		_ = metawriter.WriteError(err)
		errL.WithError(err).Fatalf("failed to GetGroups")
	}

	cfg.ChangeLog += "\n### " + time.Now().Format("2006-01-02 15:04:05")
	var releases = make([][]string, len(groups))
	for i, group := range groups {
		for j := range group.Applications {
			app := *group.Applications[j]
			releaseID, ok, err := oapi.GetLatestApplicationRelease(cfg, app)
			if err != nil {
				_ = metawriter.WriteSuccess(false)
				_ = metawriter.WriteError(err)
				errL.WithError(err).
					WithField("application name", app.Name).
					WithField("branch", app.Branch).
					Fatalf("failed to GetLatestApplicationRelease")
			}
			if !ok {
				infoL.Infof("group[%v].applications[%v], name: %s, branch: %s, latest release not found",
					i, j, app.Name, app.Branch)
				_ = metawriter.WriteWarn(fmt.Sprintf("missing group[%v].applications[%v], name: %s, branch: %s, releaseID: %s",
					i, j, app.Name, app.Branch, app.ReleaseID))
				cfg.ChangeLog += fmt.Sprintf("\n- [ ] %s %s %s", app.Name, app.Branch, app.ReleaseID)
				continue
			}
			infoL.Infof("group[%v].applications[%v], name: %s, branch: %s, releaseID: %s",
				i, j, app.Name, app.Branch, releaseID)
			releases[i] = append(releases[i], releaseID)
			cfg.ChangeLog += fmt.Sprintf("\n- [x] %s %s %s", app.Name, app.Branch, app.ReleaseID)
		}
	}
	cfg.ChangeLog += "\n"

	releaseID, err := oapi.CreateUpdateProjectRelease(cfg, releases)
	if err != nil {
		_ = metawriter.WriteSuccess(false)
		_ = metawriter.WriteError(err)
		errL.WithError(err).
			WithField("releases", releases).
			Fatalf("failed to CreateUpdateProjectRelease")
	}
	_ = metawriter.WriteSuccess(true)
	_ = metawriter.WriteKV("releaseID", releaseID)
	_ = metawriter.WriteKV("version", cfg.Version)
}
