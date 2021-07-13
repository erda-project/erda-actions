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
	"time"

	"github.com/erda-project/erda/pkg/database/sqllint/configuration"
	"github.com/erda-project/erda/pkg/database/sqllint/rules"
	"github.com/erda-project/erda/pkg/database/sqlparser/migrator"
	"github.com/erda-project/erda/pkg/envconf"
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
	Database      string   `env:"ACTION_DATABASE"`
	MigrationDir_ string   `env:"ACTION_MIGRATIONDIR"`
	NeedMySQLLint bool     `env:"ACTION_MYSQLLINT"`
	LintConfig    string   `env:"ACTION_LINT_CONFIG"`
	Modules_      []string `env:"ACTION_MODULES"`

	mysqlParameters   *migrator.DSNParameters
	sandboxParameters *migrator.DSNParameters
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

	conf.mysqlParameters = &migrator.DSNParameters{
		Database:  conf.Database,
		ParseTime: true,
		Timeout:   time.Second * 150,
	}
	conf.sandboxParameters = &migrator.DSNParameters{
		Username:  "root",
		Password:  "12345678",
		Host:      "0.0.0.0",
		Port:      3306,
		Database:  conf.Database,
		ParseTime: true,
		Timeout:   time.Second * 150,
	}
	if err := conf.retrieveMySQLParameters(); err != nil {
		logrus.Fatalf("failed to get MySQL addon DSN: %v", err)
	}

	return conf
}

// MySQLParameters returns MySQL addon's settings
func (c *Conf) MySQLParameters() *migrator.DSNParameters {
	return c.mysqlParameters
}

// SandboxParameters returns sandbox's settings
func (c *Conf) SandboxParameters() *migrator.DSNParameters {
	return c.sandboxParameters
}

// MigrationDir returns migration scripts direction like .dice/migrations or migrations
func (c *Conf) MigrationDir() string {
	return c.MigrationDir_
}

// DebugSQL returns weather to debug SQL executing
func (c *Conf) DebugSQL() bool {
	return c.PipelineDebugMode
}

func (c *Conf) Workdir() string {
	return c.WorkDir
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

func (c *Conf) retrieveMySQLParameters() error {
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

			c.mysqlParameters.Username = detail.Config.MySQLUserName
			c.mysqlParameters.Password = detail.Config.MySQLPassword
			c.mysqlParameters.Host = detail.Config.MySQLHost
			port, err := strconv.ParseUint(detail.Config.MySQLPort, 10, 32)
			if err != nil {
				return errors.Wrapf(err, "failed to parse MySQL port")
			}
			c.mysqlParameters.Port = int(port)

			return nil
		}
	}

	return errors.Errorf("mysql addon not found, applicationID: %v, workspace: %s", c.AppID, c.Workspace)
}
