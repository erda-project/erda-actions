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
package repo

import (
	"path/filepath"
	"strings"

	"github.com/erda-project/erda/pkg/parser/diceyml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/oapi"
)

type ReleasedYaml struct {
	Conf       *config.Config
	ReplaceOld string
	ReplaceNew string
	text       []byte
}

func (y *ReleasedYaml) ReadFromDiceHub(api *oapi.AccessAPI) error {
	if api == nil {
		return errors.New("AccessAPI is nil")
	}

	header := api.RequestHeader()
	header.Add("Accept", "application/x-yaml")
	data, _, err := oapi.RequestGet(api.GetDiceURL(), header)
	if err != nil {
		return err
	}

	y.text = data
	return nil
}

func (y *ReleasedYaml) Deployable() (string, error) {
	deployable, err := diceyml.NewDeployable(y.text, diceyml.WS_PROD, false)
	if err != nil {
		return "", errors.Wrap(err, "failed to NewDeployable")
	}
	return y.replaceRegistry(deployable)
}

func (y *ReleasedYaml) replaceRegistry(dice *diceyml.DiceYaml) (string, error) {
	if y.ReplaceNew == "" {
		return dice.YAML()
	}

	obj := dice.Obj()
	if y.ReplaceOld == "" {
		for name, service := range obj.Services {
			oldImage := service.Image
			if firstSlashIndex := strings.Index(service.Image, "/"); firstSlashIndex >= 0 {
				service.Image = y.ReplaceNew + service.Image[firstSlashIndex:]
			}
			logrus.WithFields(logrus.Fields{
				"service name": name,
				"old":          oldImage,
				"new":          service.Image,
			}).Infoln("replace registry")
		}
	} else {
		for name, service := range obj.Services {
			oldImage := service.Image
			service.Image = strings.ReplaceAll(service.Image, y.ReplaceOld, y.ReplaceNew)
			logrus.WithFields(logrus.Fields{
				"service name": name,
				"old":          oldImage,
				"new":          service.Image,
			}).Infoln("replace registry")
		}
	}

	out, err := yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (y *ReleasedYaml) Bucket() string {
	return y.Conf.OssBucket
}

func (y *ReleasedYaml) Local() string {
	return "dice.yml"
}

// Remote is like /archived-versions/{git-tag:v1.0.0}/releases/{repo-name:erda}/dice.yml
func (y *ReleasedYaml) Remote() string {
	return filepath.Join(y.Conf.GetOssPath(), "releases", y.Conf.GetReleaseName(), "dice.yml")
}
