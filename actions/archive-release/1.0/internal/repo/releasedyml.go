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

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/erda-project/erda/pkg/parser/diceyml"

	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/oapi"
)

type ReleasedYaml struct {
	conf       *config.Config
	replaceOld string
	replaceNew string
	obj        *diceyml.Object
}

func NewReleasedYaml(conf *config.Config) *ReleasedYaml {
	return &ReleasedYaml{conf: conf}
}

func (y *ReleasedYaml) SetReplacement(src, dst string) {
	y.replaceOld = src
	y.replaceNew = dst
}

// ReadFromDiceHub read dice.yml from dicehub and make it deployable
func (y *ReleasedYaml) ReadFromDiceHub(api *oapi.AccessAPI) (string, error) {
	if api == nil {
		return "", errors.New("AccessAPI is nil")
	}

	header := api.RequestHeader()
	header.Add("Accept", "application/x-yaml")
	data, _, err := oapi.RequestGet(api.GetDiceURL(), header)
	if err != nil {
		return "", err
	}

	return y.deployable(data)
}

// Obj returns the dice.yml structure
func (y *ReleasedYaml) Obj() *diceyml.Object {
	return y.obj
}

func (y *ReleasedYaml) Bucket() string {
	return y.conf.OssBucket
}

func (y *ReleasedYaml) Local() string {
	return "dice.yml"
}

// Remote is like /archived-versions/{git-tag:v1.0.0}/releases/{repo-name:erda}/dice.yml
func (y *ReleasedYaml) Remote() string {
	return filepath.Join(y.conf.GetOssPath(), "releases", y.conf.GetReleaseName(), "dice.yml")
}

// deployable dose
// - make make dice.yml deployable for WS_PROD;
// - patch securityContext.privileged=true to the specified service;
// - replace registry for every service's image
func (y *ReleasedYaml) deployable(text []byte) (string, error) {
	deployable, err := diceyml.NewDeployable(text, diceyml.WS_PROD, false)
	if err != nil {
		return "", errors.Wrap(err, "failed to NewDeployable")
	}

	y.obj = deployable.Obj()
	patchSecurityContextPrivileged(y.obj, y.conf.SecurityCtx...)
	return y.replaceRegistry(y.obj)
}

func (y *ReleasedYaml) replaceRegistry(obj *diceyml.Object) (string, error) {
	if y.replaceNew == "" {
		out, err := yaml.Marshal(obj)
		if err != nil {
			return "", err
		}
		return string(out), nil
	}
	if y.replaceOld == "" {
		for name, service := range obj.Services {
			oldImage := service.Image
			if firstSlashIndex := strings.Index(service.Image, "/"); firstSlashIndex >= 0 {
				service.Image = y.replaceNew + service.Image[firstSlashIndex:]
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
			service.Image = strings.ReplaceAll(service.Image, y.replaceOld, y.replaceNew)
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

func patchSecurityContextPrivileged(obj *diceyml.Object, services ...string) {
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
