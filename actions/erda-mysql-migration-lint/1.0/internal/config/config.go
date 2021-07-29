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

	"github.com/erda-project/erda/pkg/database/sqllint/configuration"
	"github.com/erda-project/erda/pkg/database/sqllint/rules"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var conf *Conf

type Conf struct {
	// basic envs
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE" required:"true"`

	PipelineDebugMode bool `env:"PIPELINE_DEBUG_MODE"`

	// action envs
	WorkDir       string `env:"ACTION_WORKDIR"`
	MigrationDir_ string `env:"ACTION_MIGRATIONDIR"`
	LintConfig    string `env:"ACTION_LINT_CONFIG"`
}

func Configuration() *Conf {
	if conf != nil {
		return conf
	}

	conf = new(Conf)
	if err := envconf.Load(conf); err != nil {
		logrus.Fatalf("failed to load configuration: %v", err)
	}

	logrus.SetLevel(logrus.InfoLevel)
	if conf.PipelineDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return conf
}

func ConfigurationString() string {
	c := Configuration()
	if c == nil {
		return ""
	}
	data, err := yaml.Marshal(*c)
	if err != nil {
		return ""
	}
	return string(data)
}

func (c *Conf) Workdir() string {
	return c.WorkDir
}

// MigrationDir returns migration scripts direction like .dice/migrations or migrations
func (c *Conf) MigrationDir() string {
	return c.MigrationDir_
}

func (c *Conf) Modules() []string {
	return nil
}

func (c *Conf) Rules() []rules.Ruler {
	configFilename := filepath.Join(c.Workdir(), c.LintConfig)
	rulesConfig, err := configuration.FromLocal(configFilename)
	if err != nil {
		logrus.WithError(err).Warnln("failed to load migration linter configuration from local, use default")
		return configuration.DefaultRulers()
	}
	rulers, err := rulesConfig.Rulers()
	if err != nil {
		logrus.WithError(err).Warnln("failed to parse migration linter from local, use default")
		return configuration.DefaultRulers()
	}
	return rulers
}
