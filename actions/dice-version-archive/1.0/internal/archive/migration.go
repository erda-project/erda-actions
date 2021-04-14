// 读取 migration 下的文件

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
