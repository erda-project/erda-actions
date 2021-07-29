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
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	configuration2 "github.com/erda-project/erda/pkg/database/sqllint/configuration"
	"github.com/erda-project/erda/pkg/database/sqllint/rules"
	"github.com/erda-project/erda/pkg/database/sqlparser/migrator"
	"github.com/erda-project/erda/pkg/envconf"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	versionPackage  = "/opt/dice-tools/versionpackage"
	versionFilename = versionPackage + "/version"
)

var configuration *Configuration

type Configuration struct {
	envs *envs
	cf   *ConfigFile
}

func Config() *Configuration {
	if configuration == nil {
		configuration = new(Configuration)
		if err := configuration.reload(); err != nil {
			log.Errorf("failed to reload configuration: %v", err)
		}

		log.Infof("%+v", *configuration.envs)
	}
	return configuration
}

func (c Configuration) MySQLParameters() *migrator.DSNParameters {
	if c.envs.MySQLUser != "" {
		return &migrator.DSNParameters{
			Username:  c.envs.MySQLUser,
			Password:  c.envs.MySQLPassword,
			Host:      c.envs.MySQLHost,
			Port:      int(c.envs.MySQLPort),
			Database:  c.envs.MySQLDiceDB,
			ParseTime: true,
			Timeout:   time.Second * 150,
		}
	}

	if c.cf == nil {
		return new(migrator.DSNParameters)
	}

	return &migrator.DSNParameters{
		Username:  c.cf.Installs.Addons.Mysql.User,
		Password:  c.cf.Installs.Addons.Mysql.Password,
		Host:      c.cf.Installs.Addons.Mysql.Host,
		Port:      c.cf.Installs.Addons.Mysql.Port,
		Database:  c.cf.Installs.Addons.Mysql.Db,
		ParseTime: true,
		Timeout:   time.Second * 150,
	}
}

func (c Configuration) SandboxParameters() *migrator.DSNParameters {
	return &migrator.DSNParameters{
		Username:  "root",
		Password:  c.envs.SandboxRootPassword,
		Host:      "0.0.0.0",
		Port:      3306,
		Database:  c.Database(),
		ParseTime: true,
		Timeout:   time.Second * 150,
	}
}

func (c Configuration) Database() string {
	if c.envs.MySQLDiceDB != "" {
		return c.envs.MySQLDiceDB
	}

	if c.cf == nil {
		return ""
	}

	return c.cf.Installs.Addons.Mysql.Db
}

// MigrationDir returns migrations scripts dir
func (c Configuration) MigrationDir() string {
	if c.envs.MigrationDir != "" {
		return c.envs.MigrationDir
	}

	data, err := ioutil.ReadFile(versionFilename)
	if err != nil {
		return ""
	}
	migrationDir := filepath.Join(versionPackage, string(data))
	return migrationDir
}

// Workdir returns workdir to join the scripts' path
func (c Configuration) Workdir() string {
	return ""
}

// DebugSQL returns whether  the process need to debug SQLs
func (c Configuration) DebugSQL() bool {
	return c.envs.DebugSQL
}

func (c Configuration) SkipMigrationLint() bool {
	return c.envs.SkipLint
}

func (c Configuration) SkipSandbox() bool {
	return c.envs.SkipSandbox
}

func (c Configuration) SkipPreMigrate() bool {
	return c.envs.SkipPreMig
}

func (c Configuration) SkipMigrate() bool {
	return c.envs.SkipMigrate
}

// Modules returns the modules for installing
func (c Configuration) Modules() []string {
	if c.envs.Modules_ == "" {
		return nil
	}
	return strings.Split(c.envs.Modules_, ",")
}

// Rules returns Erda MySQL linters
// note: hard code here
func (c Configuration) Rules() []rules.Ruler {
	return configuration2.DefaultRulers()
}

// reload reloads the envs and ${DICE_CONFIG}/config.yaml
func (c *Configuration) reload() error {
	c.envs = new(envs)
	if err := envconf.Load(c.envs); err != nil {
		return errors.Wrap(err, "failed to Load envs")
	}

	if data, err := ioutil.ReadFile(c.envs.ConfigPath); err == nil {
		c.cf = new(ConfigFile)
		_ = yaml.Unmarshal(data, c.cf) // allows err
	}

	return nil
}

type envs struct {
	ConfigPath string `env:"CONFIGPATH"` // ${DICE_CONFIG}/config.yaml

	// mysql server parameters
	MySQLUser     string `env:"MYSQL_USER"`
	MySQLPassword string `env:"MYSQL_PASSWORD"`
	MySQLHost     string `env:"MYSQL_HOST"`
	MySQLPort     uint64 `env:"MYSQL_PORT"`
	MySQLDiceDB   string `env:"MYSQL_DICE_DB"`

	// flow control parameters
	SkipLint    bool `env:"MIGRATION_SKIP_LINT"`
	SkipSandbox bool `env:"MIGRATION_SKIP_SANDBOX"`
	SkipPreMig  bool `env:"MIGRATION_SKIP_PRE_MIGRATION"`
	SkipMigrate bool `env:"MIGRATION_SKIP_MIGRATION"`

	DebugSQL bool   `env:"MIGRATION_DEBUGSQL"`
	Modules_ string `env:"MIGRATION_MODULES"`

	SandboxRootPassword string `env:"MYSQL_ROOT_PASSWORD"`

	Workdir      string `env:"WORKDIR"`
	MigrationDir string `env:"MIGRATION_DIR"`
}

// ConfigFile represents the structure of ${DICE_CONFIG}/config.yaml which
// can be read mysql configurations from .
type ConfigFile struct {
	Version  string `json:"version" yaml:"version"`
	Installs struct {
		DataDir    string `json:"data_dir" yaml:"data_dir"`
		NetdataDir string `json:"netdata_dir" yaml:"netdata_dir"`
		Addons     struct {
			Mysql struct {
				Host     string `json:"host" yaml:"host"`
				Port     int    `json:"port" yaml:"port"`
				User     string `json:"user" yaml:"user"`
				Password string `json:"password" yaml:"password"`
				Db       string `json:"db" yaml:"db"`
			} `json:"mysql" yaml:"mysql"`
		} `json:"addons" yaml:"addons"`
	} `json:"installs" yaml:"installs"`
}
