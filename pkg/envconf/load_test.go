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

package envconf_test

import (
	"os"
	"testing"

	"github.com/erda-project/erda-actions/pkg/envconf"
)

const (
	MYSQL_USERNAME = "user_cm"
	MIGRATION_MYSQL_USERNAME = "user_conf"
MYSQL_PASSWORD = "pass_cn"
	MIGRATION_MYSQL_PASSWORD = "pass_conf"
MYSQL_HOST = "host.cm.rds"
	MIGRATION_MYSQL_HOST = "127.0.0.1"
MYSQL_PORT = "3306"
	MIGRATION_MYSQL_PORT = "3306"
MYSQL_DATABASE = "erda_cm"
	MIGRATION_MYSQL_DBNAME = "erda_conf"
)

type env struct {
	MySQLUser     string `env:"MYSQL_USERNAME:MIGRATION_MYSQL_USERNAME"`
	MySQLPassword string `env:"MYSQL_PASSWORD:MIGRATION_MYSQL_PASSWORD"`
	MySQLHost     string `env:"MYSQL_HOST:MIGRATION_MYSQL_HOST"`
	MySQLPort     uint64 `env:"MYSQL_PORT:MIGRATION_MYSQL_PORT"`
	MySQLDiceDB   string `env:"MYSQL_DATABASE:MIGRATION_MYSQL_DBNAME"`
}

func initEnv1() {
	os.Setenv("MYSQL_USERNAME", MYSQL_USERNAME)
	os.Setenv("MYSQL_PASSWORD", MYSQL_PASSWORD)
	os.Setenv("MYSQL_HOST", MYSQL_HOST)
	os.Setenv("MYSQL_PORT", MYSQL_PORT)
	os.Setenv("MYSQL_DATABASE", MYSQL_DATABASE)
}

func initEnv2() {
	os.Unsetenv("MYSQL_USERNAME")
	os.Unsetenv("MYSQL_PASSWORD")
	os.Setenv("MIGRATION_MYSQL_USERNAME", MIGRATION_MYSQL_USERNAME)
	os.Setenv("MIGRATION_MYSQL_PASSWORD", MIGRATION_MYSQL_PASSWORD)
}

func initEnv3() {
	os.Setenv("MIGRATION_MYSQL_HOST", MIGRATION_MYSQL_HOST)
	os.Setenv("MIGRATION_MYSQL_DBNAME", MIGRATION_MYSQL_DBNAME)
}

func initEnv4() {
	os.Unsetenv("MYSQL_HOST")
	os.Unsetenv("MYSQL_PORT")
	os.Unsetenv("MYSQL_DATABASE")
}

func TestLoad(t *testing.T) {
	initEnv1()
	e := new(env)
	if err := envconf.Load(e); err != nil {
		t.Fatalf("failed to Load: %v", err)
	}
	t.Logf("%+v", e)
	if e.MySQLUser != MYSQL_USERNAME || e.MySQLPassword != MYSQL_PASSWORD || e.MySQLDiceDB != MYSQL_DATABASE {
		t.Fatal("initEnv1: Load env data error")
	}

	initEnv2()
	if err := envconf.Load(e); err != nil {
		t.Fatalf("failed to Load: %v", err)
	}
	t.Logf("%+v", e)
	if e.MySQLUser != MIGRATION_MYSQL_USERNAME || e.MySQLPassword != MIGRATION_MYSQL_PASSWORD || e.MySQLDiceDB != MYSQL_DATABASE {
		t.Fatal("initEnv2: Load env data error")
	}

	initEnv3()
	if err := envconf.Load(e); err != nil {
		t.Fatalf("failed to Load: %v", err)
	}
	t.Logf("%+v", e)
	if e.MySQLUser != MIGRATION_MYSQL_USERNAME || e.MySQLPassword != MIGRATION_MYSQL_PASSWORD || e.MySQLDiceDB != MYSQL_DATABASE {
		t.Fatal("initEnv2: Load env data error")
	}

	initEnv4()
	if err := envconf.Load(e); err != nil {
		t.Fatalf("failed to Load: %v", err)
	}
	t.Logf("%+v", e)
	if e.MySQLUser != MIGRATION_MYSQL_USERNAME || e.MySQLPassword != MIGRATION_MYSQL_PASSWORD || e.MySQLDiceDB != MIGRATION_MYSQL_DBNAME {
		t.Fatal("initEnv2: Load env data error")
	}
}
