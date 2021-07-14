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

package config

import (
	"path/filepath"
	"strings"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/pkg/errors"
)

const (
	Success = "success"
	Err     = "err"
)

var conf *Config

func New() (*Config, error) {
	if conf != nil {
		return conf, nil
	}

	conf = new(Config)
	if err := envconf.Load(conf); err != nil {
		return nil, errors.Wrap(err, "failed to Load envs")
	}

	return conf, nil
}

type Config struct {
	// basic envs
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE" required:"true"`

	// actions parameters
	Repos          []string `env:"ACTION_REPOS" required:"true"`
	OssEndpoint    string   `env:"ACTION_OSSENDPOINT" required:"true"`
	OssBucket      string   `env:"ACTION_OSSBUCKET" required:"true"`
	OssPath        string   `env:"ACTION_OSSPATH" required:"false"`
	OssKey         string   `env:"ACTION_OSSACCESSKEYID" required:"true"`
	OssSecret      string   `env:"ACTION_OSSACCESSKEYSECRET" required:"true"`
	OssArchivedDir string   `env:"ACTION_OSSARCHIVEDDIR" required:"true"`
	GitRef         string   `env:"ACTION_GITREF" required:"true"`

	// other parameters
	MetaFilename string `env:"METAFILE"`
}

func (c Config) GetOssPath() string {
	if c.OssPath != "" {
		return c.OssPath
	}
	if c.OssArchivedDir == "" {
		c.OssArchivedDir = "archived-versions"
	}

	version := "v" + strings.TrimPrefix(filepath.Base(c.GitRef), "v")
	return filepath.Join(c.OssArchivedDir, version, "extensions")
}
