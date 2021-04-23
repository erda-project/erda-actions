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
	"strconv"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
)

// metafile keys
const (
	MrID    = "mr_id"
	Success = "success"
	Err     = "err"
	Warn    = "warn"
	Step    = "step"
)

const (
	DiceYmlPathFromSrcRepo             = "dice.yml"
	MigrationPathFromDstRepoVersionDir = "sqls"
)

var c *config

func init() {
	initLog()
}

type config struct {
	// basic envs
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE" required:"true"`

	// pipeline parameters
	PipelineDebugMode bool   `env:"PIPELINE_DEBUG_MODE"`
	PipelineID        string `env:"PIPELINE_ID"`
	PipelineTaskLogID string `env:"PIPELINE_TASK_LOG_ID"`
	PipelineTaskID    string `env:"PIPELINE_TASK_ID"`

	// action parameters
	Workdir                   string               `env:"ACTION_WORKDIR"`
	MigrationsPathFromSrcRepo string               `env:"ACTION_MIGRATIONS_DIR"`
	Dst                       RepoInfo             `env:"ACTION_DST"`
	MRProcessor               uint64               `env:"ACTION_MR_PROCESSOR"`
	Registry                  *RegistryReplacement `env:"ACTION_REGISTRY_REPLACEMENT"`
	ReleaseID                 string               `env:"ACTION_RELEASEID"`

	// other parameters
	MetaFilename string `env:"METAFILE"`
}

type RepoInfo struct {
	RepoName string `json:"repoName"`
	Branch   string `json:"branch"`
	SnapName string `json:"snapName"`
}

type RegistryReplacement struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func configuration() *config {
	if c == nil {
		c = new(config)
		if err := envconf.Load(c); err != nil {
			logrus.Errorf("failed to load configuration, err: %v", err)
		}
		if c.Dst.SnapName == "" {
			c.Dst.SnapName = c.AppName
		}
		if c.Dst.Branch == "" {
			c.Dst.Branch = "master"
		}
	}

	return c
}

func OrgID() uint64 {
	return configuration().OrgID
}

func OpenapiToken() string {
	return configuration().CiOpenapiToken
}

func OpenapiPrefix() string {
	return configuration().DiceOpenapiPrefix
}

func Workdir() string {
	return configuration().Workdir
}

func MigrationsPathFromSrcRepoRoot() string {
	if p := configuration().MigrationsPathFromSrcRepo; p != "" {
		return p
	}
	return ".dice/migrations"
}

func PipelineID() string {
	return configuration().PipelineID
}

func ProjectID() int64 {
	return configuration().ProjectID
}

func ProjectName() string {
	return configuration().ProjectName
}

func ApplicationID() uint64 {
	return configuration().AppID
}

func ApplicationName() string {
	return configuration().AppName
}

func DstApplicationName() string {
	if configuration().Dst.RepoName == "" {
		return "version"
	}
	return configuration().Dst.RepoName
}

func DstRepoRefBranch() string {
	return configuration().Dst.Branch
}

func DstRepoBranch() string {
	return "feature/pipeline-" + configuration().PipelineID + "-" + configuration().PipelineTaskID
}

func MRProcessor() string {
	return strconv.FormatUint(configuration().MRProcessor, 10)
}

func Replacement() *RegistryReplacement {
	return configuration().Registry
}

// DiceYmlPathFromDstRepoVersionDir e.g. releases/erda/dice.yml, "erda" in the path defaults current application name
func DiceYmlPathFromDstRepoVersionDir() string {
	if configuration().Dst.SnapName == "" {
		return filepath.Join("releases", configuration().AppName, "dice.yml")
	}
	return filepath.Join("releases", configuration().Dst.SnapName, "dice.yml")
}

func ReleaseID() string {
	return configuration().ReleaseID
}

func initLog() {
	log.Init()
	if configuration().PipelineDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
