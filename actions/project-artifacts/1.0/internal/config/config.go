// Copyright (c) 2022 Terminus, Inc.
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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	"github.com/erda-project/erda/pkg/envconf"
)

var (
	c                  *Config
	defaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)

type Config struct {
	// platform envs
	OrgID     uint64 `env:"DICE_ORG_ID" required:"true"`
	OapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	OapiHost  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectID int64  `env:"DICE_PROJECT_ID" required:"true"`
	UserID    string `env:"DICE_USER_ID" required:"true"`

	// action envs
	Version   string `env:"ACTION_VERSION"`
	ChangeLog string `env:"ACTION_CHANGELOG"`
	Groups    string `env:"ACTION_GROUPS"`
	Tz        string `env:"ACTION_TZ"`
}

func (cfg Config) Host() string {
	for i := 0; i < len(cfg.OapiHost)-1; i++ {
		if cfg.OapiHost[i:i+2] == "//" {
			return cfg.OapiHost[i+2:]
		}
	}
	return cfg.OapiHost
}

func (cfg Config) GetGroups() ([]*Group, error) {
	var groups []*Group
	if err := yaml.Unmarshal([]byte(cfg.Groups), &groups); err != nil {
		return nil, errors.Wrapf(err, "failed to Unmarshal Groups: %s", cfg.Groups)
	}
	if len(groups) == 0 {
		return nil, errors.Errorf("no group in params: %s", cfg.Groups)
	}
	for i := range groups {
		if len(groups[i].Applications) == 0 {
			return nil, errors.Errorf("no application in group[%v]: %s", i, cfg.Groups)
		}
	}
	return groups, nil
}

func (cfg Config) Print() {
	_, _ = fmt.Fprintf(os.Stdout, "\t    host: %s\n", cfg.Host())
	_, _ = fmt.Fprintf(os.Stdout, "\t  verson: %s\n", cfg.Version)
	_, _ = fmt.Fprintf(os.Stdout, "\tchangeLog: %s\n", cfg.ChangeLog)
	groups, err := cfg.GetGroups()
	if err != nil {
		return
	}
	for i := range groups {
		_, _ = fmt.Fprintf(os.Stdout, "\tgroup[%v]\n", i)
		for j := range groups[i].Applications {
			_, _ = fmt.Fprintf(os.Stdout, "\t\tapplication[%v]\n", j)
			_, _ = fmt.Fprintf(os.Stdout, "\t\t       name: %s\n", groups[i].Applications[j].Name)
			_, _ = fmt.Fprintf(os.Stdout, "\t\t     branch: %s\n", groups[i].Applications[j].Branch)
			_, _ = fmt.Fprintf(os.Stdout, "\t\t  releaseID: %s\n", groups[i].Applications[j].ReleaseID)
		}
	}
}

type Group struct {
	Applications []*Application
}

type Application struct {
	Name      string `json:"name" yaml:"name"`
	Branch    string `json:"branch" yaml:"branch"`
	ReleaseID string `json:"releaseID" yaml:"releaseID"`
}

func GetConfig() (*Config, error) {
	if c != nil {
		return c, nil
	}
	c = new(Config)
	if err := envconf.Load(c); err != nil {
		return nil, errors.Wrap(err, "failed to Load params")
	}
	if c.Tz == "" {
		c.Tz = "Asia/Shanghai"
	}
	if !strings.Contains(c.Version, "+") {
		location, err := time.LoadLocation(c.Tz)
		if err != nil {
			location = defaultLocation
		}
		c.Version += "+" + time.Now().In(location).Format("20060102150405")
	}
	return c, nil
}
