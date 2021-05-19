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

package migration

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/sqllint/configuration"
	"github.com/erda-project/erda/pkg/sqllint/rules"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	addonListURI       = "/api/addons?type=project&value=%s"
	addonDetailURI     = "/api/addons/%s"
	addonReferencesURI = "/api/addons/%s/actions/references"
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
	WorkDir       string   `env:"ACTION_WORKDIR"`
	Database_     string   `env:"ACTION_DATABASE"`
	MigrationDir_ string   `env:"ACTION_MIGRATIONDIR"`
	NeedMySQLLint bool     `env:"ACTION_MYSQLLINT"`
	LintConfig    string   `env:"ACTION_LINT_CONFIG"`
	Modules_      []string `env:"ACTION_MODULES"`

	MetaFilename_ string `env:"METAFILE"`

	dsn string
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

	if err := conf.retrieveDSN(); err != nil {
		logrus.Fatalf("failed to get MySQL addon DSN: %v", err)
	}

	return conf
}

// DSN gets MySQL DSN
func (c *Conf) DSN() string {
	return c.dsn
}

// SandboxDSN gets sandbox DSN
func (c *Conf) SandboxDSN() string {
	return "root:12345678@(localhost:3306)/"
}

// MigrationDir gets migration scripts direction like .dice/migrations or migrations
func (c *Conf) MigrationDir() string {
	return c.MigrationDir_
}

// AppVersion gets application version
func (c *Conf) AppVersion() string {
	return ""
}

// BaseVersion gets base version
func (c *Conf) BaseVersion() string {
	return ""
}

// DebugSQL gets weather to debug SQL executing
func (c *Conf) DebugSQL() bool {
	return c.PipelineDebugMode
}

func (c *Conf) Database() string {
	return c.Database_
}

func (c *Conf) Workdir() string {
	return c.WorkDir
}

func (c *Conf) MetaFilename() string {
	return c.MetaFilename_
}

func (c *Conf) NeedErdaMySQLLint() bool {
	return c.NeedMySQLLint
}

func (c *Conf) Modules() []string {
	var modules []string
	for _, v := range c.Modules_ {
		ss := strings.Split(v, ",")
		for _, vv := range ss {
			modules = append(modules, strings.TrimSpace(vv))
		}
	}
	return modules
}

func (c *Conf) Rules() []rules.Ruler {
	configFilename := filepath.Join(c.Workdir(), c.LintConfig)
	rulesConfig, err := configuration.FromLocal(configFilename)
	if err != nil {
		logrus.Warnln("failed to load migration linter configuration from local, use default")
		return configuration.DefaultRulers()
	}
	rulers, err := rulesConfig.Rulers()
	if err != nil {
		logrus.Warnln("failed to parse migration linter from local, use default")
		return configuration.DefaultRulers()
	}
	return rulers
}

func (c *Conf) retrieveDSN() error {
	// 查找项目下所有的 addon 实例
	url := c.DiceOpenapiPrefix + fmt.Sprintf(addonListURI, strconv.FormatUint(uint64(c.ProjectID), 10))
	header := map[string][]string{"authorization": {c.CiOpenapiToken}}
	list, err := getAddonList(url, header)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.Errorf("there is no addon in the project, projectID: %v", c.ProjectID)
	}

	// filter mysql with the workspace
	var mysqlAddons []GetAddonsListResponseDataEle
	for _, addon := range list {
		if strings.EqualFold(addon.AddonName, "mysql") && strings.EqualFold(addon.Workspace, c.Workspace) {
			mysqlAddons = append(mysqlAddons, addon)
		}
	}
	if len(mysqlAddons) == 0 {
		return errors.Errorf("there is no MySQL addon on the current workspace %s", c.Workspace)
	}

	for _, addon := range mysqlAddons {
		url := c.DiceOpenapiPrefix + fmt.Sprintf(addonReferencesURI, addon.InstanceID)
		references, err := getAddonReferences(url, header)
		if err != nil {
			return err
		}

		for _, ref := range references {
			if ref.ApplicationID != c.AppID {
				continue
			}

			url := c.DiceOpenapiPrefix + fmt.Sprintf(addonDetailURI, addon.InstanceID)
			detail, err := getAddonDetail(url, header)
			if err != nil {
				return err
			}

			c.dsn = fmt.Sprintf("%s:%s@(%s:%s)/",
				detail.Config.MySQLUserName,
				detail.Config.MySQLPassword,
				detail.Config.MySQLHost,
				detail.Config.MySQLPort)

			return nil
		}
	}

	return errors.Errorf("mysql addon not found, applicationID: %v, workspace: %s", c.AppID, c.Workspace)
}
