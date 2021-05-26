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

package workdir

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// New returns a Workdir
func New(ctx context.Context, workdir string) *Workdir {
	return &Workdir{
		Workdir: workdir,
		ctx:     ctx,
	}
}

// Workdir is an object as a repo which contains extensions (actions or addons)
type Workdir struct {
	Workdir  string
	Versions []*Version

	versionPaths []string
	ctx          context.Context
}

// LoadExtensions loads all extensions from the repo (contains all versions below)
func (repo *Workdir) LoadExtensions() error {
	repo.locate(repo.Workdir)
	for _, dirname := range repo.versionPaths {
		version, err := NewVersion(repo.ctx, dirname)
		if err != nil {
			return errors.Wrap(err, "failed to NewVersion")
		}
		repo.Versions = append(repo.Versions, version)
	}

	if len(repo.Versions) == 0 {
		return errors.Errorf("there is no extension in the repo %s, please ensure this is an extension repo")
	}

	return nil
}

// Push pushes all extensions from the repo (contains all versions below)
func (repo *Workdir) Push() error {
	for _, ext := range repo.Versions {
		logrus.Infof("push extension to Erda from %s", ext.Dirname)
		if err := ext.Push(); err != nil {
			return errors.Wrapf(err, "failed to Push: %s", ext.Dirname)
		}
	}

	return nil
}

func (repo *Workdir) locate(dirname string) {
	infos, ok := isThereSpecFile(dirname)
	if ok {
		repo.versionPaths = append(repo.versionPaths, dirname)
		return
	}

	for _, cur := range infos {
		repo.locate(filepath.Join(dirname, cur.Name()))
	}
}

func isThereSpecFile(dirname string) ([]os.FileInfo, bool) {
	var dirs []os.FileInfo
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, false
	}
	for _, file := range infos {
		if file.IsDir() {
			dirs = append(dirs, file)
			continue
		}
		if strings.EqualFold(file.Name(), "spec.yml") || strings.EqualFold(file.Name(), "spec.yaml") {
			return nil, true
		}
	}
	return dirs, false
}
