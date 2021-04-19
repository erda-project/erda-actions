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
package archive

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Script struct {
	NameFromService string
	Content         []byte

	filename string
}

func ReadScripts(workdir, migdir string) ([]*Script, error) {
	services, err := ioutil.ReadDir(filepath.Join(workdir, migdir))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ReadDir %s", filepath.Join(workdir, migdir))
	}

	var scripts []*Script
	for _, service := range services {
		logrus.Debugln("service name:", service.Name())

		if !service.IsDir() {
			continue
		}

		files, err := ioutil.ReadDir(filepath.Join(workdir, migdir, service.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to ReadDir %s", filepath.Join(workdir, migdir, service.Name()))
		}

		for _, file := range files {
			logrus.Debugln("\tfile name:", file.Name())

			if file.IsDir() {
				continue
			}

			var script Script
			script.filename = filepath.Join(workdir, migdir, service.Name(), file.Name())
			script.NameFromService = filepath.Join(service.Name(), file.Name())
			script.Content, err = ioutil.ReadFile(script.filename)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to ReadFile %s", script.filename)
			}

			scripts = append(scripts, &script)
		}
	}

	return scripts, nil
}
