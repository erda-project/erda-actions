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

// read dice.yml from erda repo
package archive

import (
	"io/ioutil"
	"strings"

	"github.com/erda-project/erda/pkg/parser/diceyml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/erda-project/erda-actions/actions/dice-version-archive/1.0/internal/config"
)

type DiceYaml struct {
	text []byte
}

func (y *DiceYaml) Read(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	y.text = data
	return nil
}

func (y *DiceYaml) Deployable() (string, error) {
	deployable, err := diceyml.NewDeployable(y.text, diceyml.WS_PROD, false)
	if err != nil {
		return "", err
	}
	return replaceRegistry(deployable)
}

func replaceRegistry(dice *diceyml.DiceYaml) (string, error) {
	replacement := config.Replacement()
	if replacement == nil || replacement.New == "" {
		return dice.YAML()
	}

	logrus.Debugln("replace registry")
	obj := dice.Obj()
	if replacement.Old == "" {
		for name, service := range obj.Services {
			oldImage := service.Image
			if firstSlashIndex := strings.Index(service.Image, "/"); firstSlashIndex >= 0 {
				service.Image = replacement.New + service.Image[firstSlashIndex:]
			}
			logrus.Debugf("service name: %s, old iamge: %s, new iamge: %s", name, oldImage, service.Image)
		}
	} else {
		for name, service := range obj.Services {
			oldImage := service.Image
			service.Image = strings.ReplaceAll(service.Image, replacement.Old, replacement.New)
			logrus.Debugf("service name: %s, old image: %s, new image: %s", name, oldImage, service.Image)
		}
	}

	out, err := yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
