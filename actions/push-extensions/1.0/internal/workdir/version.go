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
	"path/filepath"
	"strings"

	"github.com/erda-project/erda/apistructs"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-actions/actions/push-extensions/1.0/internal/client"
	"github.com/erda-project/erda-actions/actions/push-extensions/1.0/internal/config"
)

const (
	CtxKeyConfig = "config"
	CtxKeyClient = "client"
)

// NewVersion returns a Version
// there should be a *config.Config and *client.Client in ctx,
// dirname is the directory path from root to extension version directory.
func NewVersion(ctx context.Context, dirname string) (*Version, error) {
	conf, ok := ctx.Value(CtxKeyConfig).(*config.Config)
	if !ok {
		return nil, errors.New("no config in context")
	}
	cli, ok := ctx.Value(CtxKeyClient).(*client.Client)
	if !ok {
		return nil, errors.New("no push extension client in context")
	}

	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadDir")
	}

	var version = Version{
		Name:           filepath.Base(dirname),
		Dirname:        dirname,
		Spec:           new(apistructs.Spec),
		SpecContent:    nil,
		DiceContent:    nil,
		ReadmeContent:  nil,
		SwaggerContent: nil,
		conf:           conf,
		client:         cli,
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		switch {
		case strings.EqualFold(fileInfo.Name(), "spec.yml") || strings.EqualFold(fileInfo.Name(), "spec.yaml"):
			version.SpecContent, err = ioutil.ReadFile(filepath.Join(dirname, fileInfo.Name()))
			if err != nil {
				return nil, errors.Wrap(err, "failed to ReadFile")
			}
			if err = yaml.Unmarshal(version.SpecContent, version.Spec); err != nil {
				return nil, errors.Wrap(err, "failed to parse "+fileInfo.Name())
			}

		case strings.EqualFold(fileInfo.Name(), "dice.yml") || strings.EqualFold(fileInfo.Name(), "dice.yaml"):
			version.DiceContent, _ = ioutil.ReadFile(filepath.Join(dirname, fileInfo.Name()))

		case strings.EqualFold(fileInfo.Name(), "readme.md") || strings.EqualFold(fileInfo.Name(), "readme.markdown"):
			version.ReadmeContent, _ = ioutil.ReadFile(filepath.Join(dirname, fileInfo.Name()))

		case strings.EqualFold(fileInfo.Name(), "swagger.json") || strings.EqualFold(fileInfo.Name(), "swagger.yml") ||
			strings.EqualFold(fileInfo.Name(), "swagger.yaml"):
			version.SwaggerContent, _ = ioutil.ReadFile(filepath.Join(dirname, fileInfo.Name()))
		}
	}

	if version.Spec == nil || len(version.SpecContent) == 0 {
		return nil, errors.Errorf("spec file not found in %s", dirname)
	}

	// replace registry
	if conf.Registry != "" && len(version.DiceContent) > 0 {
		content, _, err := replaceDiceRegistry(version.DiceContent, version.Spec.Type, conf.Registry)
		if err != nil {
			return nil, errors.Errorf("failed to replace registry in dice.yml, dirname: %s", dirname)
		}
		version.DiceContent = content
	}

	return &version, nil
}

// Version is a version of an Extension
type Version struct {
	Name    string
	Dirname string

	Spec          *apistructs.Spec // structure of spec.yml
	SpecContent   []byte           // content of spec.yml
	DiceContent   []byte           // content of dice.yml
	ReadmeContent []byte           // content of readme.md

	SwaggerContent []byte // content of swagger.yml
	conf           *config.Config
	client         *client.Client
}

func (v *Version) Push() error {
	var payload = apistructs.ExtensionVersionCreateRequest{
		Name:        v.Spec.Name,
		Version:     v.Spec.Version,
		SpecYml:     string(v.SpecContent),
		DiceYml:     string(v.DiceContent),
		SwaggerYml:  string(v.SwaggerContent),
		Readme:      string(v.ReadmeContent),
		Public:      v.Spec.Public,
		ForceUpdate: true,
		All:         true,
		IsDefault:   v.Spec.IsDefault,
	}
	return v.client.Push(&payload)
}

func replaceDiceRegistry(content []byte, typ string, dstRegistry string) ([]byte, map[string]string, error) {
	var (
		diceData   = make(map[string]interface{})
		pushImages = make(map[string]string)
	)
	if err := yaml.Unmarshal(content, &diceData); err != nil {
		return nil, nil, err
	}

	var (
		jobs = make(map[string]interface{})
		ok   bool
	)
	switch typ {
	case Actions:
		if jobs, ok = diceData["jobs"].(map[string]interface{}); !ok {
			return nil, nil, errors.New("failed to parse jobs as map in dice.yml file")
		}
	case Addons:
		if jobs, ok = diceData["services"].(map[string]interface{}); !ok {
			return content, pushImages, nil
		}
	default:
		return nil, nil, errors.Errorf("invalid extension type: %s", typ)
	}

	for name, v := range jobs {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, nil, errors.Errorf("failed to parse jobs.%s value as map in dice.yml file", name)
		}
		image, ok := cfg["image"]
		if !ok {
			continue
		}
		imageStr, ok := image.(string)
		if !ok {
			return nil, nil, errors.Errorf("failed to parse jobs.%s.image value as string in dice.yml file", name)
		}
		if image == "" {
			return nil, nil, errors.Errorf("invalid jobs.%s.image value", name)
		}
		newImage := newDockerImage(imageStr, dstRegistry)
		cfg["image"] = newImage
		pushImages[name] = newImage
	}

	newContent, err := yaml.Marshal(diceData)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to yaml.Marshal dice")
	}

	return newContent, pushImages, nil
}

func newDockerImage(oldImage, newRegistry string) string {
	if newRegistry == "" {
		return oldImage
	}
	if index := strings.Index(oldImage, "/"); index > 0 {
		return newRegistry + oldImage[index:]
	}
	return oldImage
}
