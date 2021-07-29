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

// defines configurations
package config

import (
	"path/filepath"
	"strings"

	"github.com/erda-project/erda/pkg/envconf"
)

// metafile keys
const (
	Success = "success"
	Err     = "error"
	Warn    = "warn"
)

var c *Config

func New() (*Config, error) {
	if c != nil {
		return c, nil
	}

	c = new(Config)
	if err := envconf.Load(c); err != nil {
		return nil, err
	}

	return c, nil
}

type Config struct {
	// basic envs
	OrgID        uint64 `env:"DICE_ORG_ID" required:"true"`
	OpenapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	OpenapiHost  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName  string `env:"DICE_PROJECT_NAME" required:"true"`
	AppName      string `env:"DICE_APPLICATION_NAME" required:"true"`
	ProjectID    int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID        uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace    string `env:"DICE_WORKSPACE" required:"true"`

	// pipeline parameters
	PipelineDebugMode bool   `env:"PIPELINE_DEBUG_MODE"`
	PipelineID        string `env:"PIPELINE_ID"`
	PipelineTaskLogID string `env:"PIPELINE_TASK_LOG_ID"`
	PipelineTaskID    string `env:"PIPELINE_TASK_ID"`

	// action parameters
	Workdir        string               `env:"ACTION_WORKDIR"`
	MigDir         string               `env:"ACTION_MIGRATIONSDIR"`
	LintConfig     string               `env:"ACTION_LINT_CONFIG"`
	Registry       *RegistryReplacement `env:"ACTION_REGISTRY_REPLACEMENT"`
	ReleaseID      string               `env:"ACTION_RELEASEID"`
	OssEndpoint    string               `env:"ACTION_OSSENDPOINT" required:"true"`
	OssBucket      string               `env:"ACTION_OSSBUCKET" required:"true"`
	OssPath        string               `env:"ACTION_OSSPATH" required:"false"`
	OssKey         string               `env:"ACTION_OSSACCESSKEYID" required:"true"`
	OssSecret      string               `env:"ACTION_OSSACCESSKEYSECRET" required:"true"`
	OssArchivedDir string               `env:"ACTION_OSSARCHIVEDDIR" required:"true"`
	GitRef         string               `env:"ACTION_GITREF" required:"true"`
	ReleaseName    string               `env:"ACTION_RELEASENAME"`

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
	return filepath.Join(c.OssArchivedDir, version)
}

func (c Config) GetReleaseName() string {
	if c.ReleaseName != "" {
		return c.ReleaseName
	}
	return filepath.Base(c.Workdir)
}

func (c Config) GitRefDir() GitRefDir {
	return GitRefDir{
		ossBucket: c.OssBucket,
		remote:    c.GetOssPath(),
	}
}

type RegistryReplacement struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type GitRefDir struct {
	ossBucket, remote string
}

func (c GitRefDir) Bucket() string {
	return c.ossBucket
}

func (c GitRefDir) Remote() string {
	return c.remote
}

func (c GitRefDir) Local() string {
	return "nothing"
}
