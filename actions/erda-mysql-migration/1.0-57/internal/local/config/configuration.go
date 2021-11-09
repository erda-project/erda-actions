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
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/erda-project/erda-actions/pkg/envconf"
	configuration2 "github.com/erda-project/erda/pkg/database/sqllint/configuration"
	"github.com/erda-project/erda/pkg/database/sqllint/rules"
	"github.com/erda-project/erda/pkg/database/sqlparser/migrator"
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

// MySQLParameters 返回要应用数据库迁移的 MySQL Server 的 DSN 信息.
// 如果从环境变量中读取的 DSN 信息有效，则返回环境变量中的 DSN 信息;
// 否则返回从配置文件中读取到的 MySQL Addon 的 DSN 信息.
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
	if c.ExternalSandbox() {
		return &migrator.DSNParameters{
			Username:  c.envs.SandboxUsername,
			Password:  c.envs.SandboxPassword,
			Host:      c.envs.SandboxHost,
			Port:      c.envs.SandboxPort,
			Database:  c.Database(),
			ParseTime: true,
			Timeout:   time.Second * 150,
		}
	}
	return &migrator.DSNParameters{
		Username:  "root",
		Password:  c.envs.SandboxInnerPassword,
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
	data = bytes.TrimRight(data, "\n")
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

func (c *Configuration) RetryTimeout() uint64 {
	return c.envs.RetryTimout
}

func (c *Configuration) SQLCollectorDir() string {
	return "/log"
}

// Rules returns Erda MySQL linters
// note: hard code here
func (c Configuration) Rules() []rules.Ruler {
	return configuration2.DefaultRulers()
}

func (c Configuration) ExternalSandbox() bool {
	return c.envs.ExternalSandbox
}

// reload reloads the envs and ${DICE_CONFIG}/config.yaml
func (c *Configuration) reload() error {
	c.envs = new(envs)
	if err := envconf.Load(c.envs); err != nil {
		return errors.Wrap(err, "failed to Load envs")
	}

	c.envs.ConfigPath = os.Getenv("ConfigPath")
	if data, err := ioutil.ReadFile(c.envs.ConfigPath); err == nil {
		c.cf = new(ConfigFile)
		_ = yaml.Unmarshal(data, c.cf) // allows err
	}

	return nil
}

type envs struct {
	ConfigPath string `env:"CONFIGPATH"` // ${DICE_CONFIG}/config.yaml

	// mysql server parameters
	// Preferred to use env from ConfigMap dice-addons-info
	MySQLUser     string `env:"MIGRATION_MYSQL_USERNAME:MYSQL_USERNAME"`
	MySQLPassword string `env:"MIGRATION_MYSQL_PASSWORD:MYSQL_PASSWORD"`
	MySQLHost     string `env:"MIGRATION_MYSQL_HOST:MYSQL_HOST"`
	MySQLPort     uint64 `env:"MIGRATION_MYSQL_PORT:MYSQL_PORT"`
	MySQLDiceDB   string `env:"MIGRATION_MYSQL_DBNAME:MYSQL_DATABASE"`

	// flow control parameters
	SkipLint    bool `env:"MIGRATION_SKIP_LINT"`
	SkipSandbox bool `env:"MIGRATION_SKIP_SANDBOX"`
	SkipPreMig  bool `env:"MIGRATION_SKIP_PRE_MIGRATION"`
	SkipMigrate bool `env:"MIGRATION_SKIP_MIGRATION"`

	DebugSQL bool   `env:"MIGRATION_DEBUGSQL"`
	Modules_ string `env:"MIGRATION_MODULES"`

	// sandbox envs
	ExternalSandbox bool   `env:"MIGRATION_EXTERNAL_SANDBOX"`
	SandboxHost     string `env:"MIGRATION_SANDBOX_HOST"`
	SandboxPort     int    `env:"MIGRATION_SANDBOX_PORT"`
	SandboxUsername string `env:"MIGRATION_SANDBOX_USERNAME"`
	SandboxPassword string `env:"MIGRATION_SANDBOX_PASSWORD"`

	// RetryTimout is the max duration for connection to the MySQL Server and the Sandbox
	RetryTimout uint64 `env:"MIGRATION_RETRY_TIMEOUT"`

	SandboxRootPassword  string `env:"MYSQL_ROOT_PASSWORD"`
	SandboxInnerPassword string `env:"SANDBOX_INNER_PASSWORD"`

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
