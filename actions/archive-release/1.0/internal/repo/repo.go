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

package repo

import (
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/oapi"
)

func New(conf *config.Config) (*Repo, error) {
	api := oapi.NewAccessAPI(
		conf.OpenapiHost,
		conf.OpenapiToken,
		strconv.FormatUint(conf.OrgID, 10),
		conf.ReleaseID,
	)

	// read released yml from dicehub, make it deployable, and write it local
	releasedYml := new(ReleasedYaml)
	releasedYml.Conf = conf
	if replacement := conf.Registry; replacement != nil {
		releasedYml.ReplaceOld = replacement.Old
		releasedYml.ReplaceNew = replacement.New
	}
	if err := releasedYml.ReadFromDiceHub(api); err != nil {
		return nil, errors.Wrap(err, "failed to ReadFromDiceHub")
	}
	deployableContent, err := releasedYml.Deployable()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make released yml Deployable")
	}
	if err = writeReleasedYml(releasedYml.Local(), deployableContent); err != nil {
		return nil, errors.Wrap(err, "failed to write released yml to local")
	}

	// loads migrations (lint config and scripts)
	scripts, err := ReadScripts(conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadScripts")
	}
	return &Repo{
		ReleaseYml: releasedYml,
		LintConfig: &LintConfig{conf: conf},
		Scripts:    scripts,
	}, nil
}

type Repo struct {
	ReleaseYml *ReleasedYaml
	LintConfig *LintConfig
	Scripts    []*Script
}

func writeReleasedYml(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to OpenFile for writing, filename: %s", filename)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return errors.Wrapf(err, "failed to WriteString to %s", filename)
	}

	return nil
}
