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
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/sirupsen/logrus"
)

var c *config

func init() {
	initLog()
}

type config struct {

	// pipeline parameters
	PipelineDebugMode bool `env:"PIPELINE_DEBUG_MODE"`

	// erda version
	ErdaVersion string `env:"ACTION_ERDA_VERSION"`

	// repo path info in action job volume
	RepoErdaRelease string `env:"ACTION_REPO_ERDA_RELEASE"`
	RepoVersion     string `env:"ACTION_REPO_VERSION"`

	// release type env
	ReleaseType string `env:"ACTION_RELEASE_TYPE"`

	// oss auth info
	OSS *pkg.Oss `env:"ACTION_OSS"`

	// github auth info
	Git *git `env:"ACTION_GIT"`

	// other parameters
	MetaFilename string `env:"METAFILE"`
}

type git struct {
	Account string `json:"account"`
	Token   string `json:"token"`
}

func configuration() *config {
	if c == nil {
		c = new(config)
		if err := envconf.Load(c); err != nil {
			logrus.Errorf("failed to load configuration, err: %v", err)
		}
	}

	return c
}

func ErdaVersion() string {
	return configuration().ErdaVersion
}

func RepoErdaRelease() string {
	return configuration().RepoErdaRelease
}

func ReleaseType() string {

	releaseType := configuration().ReleaseType

	if releaseType != pkg.ReleaseCompletely && releaseType != pkg.ReleaseTools {
		configuration().ReleaseType = pkg.ReleaseCommon
	}

	return configuration().ReleaseType
}

func RepoVersion() string {
	return configuration().RepoVersion
}

func OssInfo() *pkg.Oss {

	oss := configuration().OSS

	if oss.OssEndPoint == "" {
		oss.OssEndPoint = "oss.aliyuncs.com"
	}

	return oss
}

func GitInfo() *git {
	return configuration().Git
}

func MetaFile() string {
	return configuration().MetaFilename
}

func initLog() {
	log.Init()
	if configuration().PipelineDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
