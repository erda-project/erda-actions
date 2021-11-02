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
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

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

	obj := deployable.Obj()
	PatchSecurityContextPrivileged(obj, y.Conf.SecurityCtx...)
	return y.replaceRegistry(obj)
}

func (y *ReleasedYaml) replaceRegistry(obj *diceyml.Object) (string, error) {
	if y.ReplaceNew == "" {
		out, err := yaml.Marshal(obj)
		if err != nil {
			return "", err
		}
		return string(out), nil
	}
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

func PatchSecurityContextPrivileged(obj *diceyml.Object, services ...string) {
	b := true
	for _, serviceName := range services {
		if service := obj.Services[serviceName]; service != nil {
			service.K8SSnippet = &diceyml.K8SSnippet{
				Container: &diceyml.ContainerSnippet{
					SecurityContext: &v1.SecurityContext{
						Privileged: &b,
					},
				},
			}
		}
	}
}
