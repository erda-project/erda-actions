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
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/erda-project/erda-actions/pkg/envconf"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

// metafile keys
const (
	Name              = "name"
	Tag               = "tag"
	Configs           = "configs"
	InstanceID        = "instanceID"
	RoutingInstanceID = "routingInstanceID"
)

var (
	c     *Config
	count int
)

func Get() *Config {
	if c != nil {
		return c
	}
	c = new(Config)
	if err := envconf.Load(c); err != nil {
		panic(errors.Wrap(err, "failed to Load config"))
	}
	if err := c.init(); err != nil {
		_ = metawriter.WriteError(errors.Wrap(err, "failed to init config"))
		logrus.Fatalf("failed to init config: %v", err)
	}
	return c
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
	Name        string            `env:"ACTION_NAME" required:"true"`
	Tag         string            `env:"ACTION_TAG"`
	Configs     map[string]string `env:"ACTION_CONFIGS"`
	ConfigsFrom string            `env:"ACTION_CONFIGSFROM"`

	realConfigs map[string]string
}

func (c *Config) init() error {
	c.realConfigs = make(map[string]string)
	var cf = configsFrom{
		Name:    "",
		Tag:     "",
		Default: make(map[string]string),
		Dev:     make(map[string]string),
		Test:    make(map[string]string),
		Staging: make(map[string]string),
		Prod:    make(map[string]string),
	}
	if c.ConfigsFrom != "" {
		file, err := ioutil.ReadFile(c.ConfigsFrom)
		if err != nil {
			_ = metawriter.WriteError(fmt.Sprintf("failed to parse configs from file: %v", err))
			logrus.WithError(err).
				WithField("configsFrom", c.ConfigsFrom).
				Fatalln("failed to ReadFile")
		}
		if err = yaml.Unmarshal(file, &cf); err != nil {
			_ = metawriter.WriteError(fmt.Sprintf("failed to parse configs from file: %v", err))
			logrus.WithError(err).
				WithField("configsFrom", c.ConfigsFrom).
				Fatalln("failed to parse config from file")
		}
	}

	if c.Name == "" {
		c.Name = cf.Name
	}
	if c.Tag == "" {
		c.Tag = cf.Tag
	}

	for k, v := range cf.Default {
		c.realConfigs[k] = v
	}
	var envConfig = make(map[string]string)
	switch strings.ToLower(c.Workspace) {
	case "prod":
		envConfig = cf.Prod
	case "staging":
		envConfig = cf.Staging
	case "test":
		envConfig = cf.Test
	case "dev":
		envConfig = cf.Dev
	}
	for k, v := range envConfig {
		c.realConfigs[k] = v
	}
	if len(c.Configs) > 0 {
		for k, v := range c.Configs {
			c.realConfigs[k] = v
		}
	}

	return Interpolate(c.realConfigs)
}

func (c *Config) GetConfigs() map[string]string {
	return c.realConfigs
}

type configsFrom struct {
	Name    string            `yaml:"name"`
	Tag     string            `yaml:"tag"`
	Default map[string]string `yaml:"default"`
	Dev     map[string]string `yaml:"dev"`
	Test    map[string]string `yaml:"test"`
	Staging map[string]string `yaml:"staging"`
	Prod    map[string]string `yaml:"prod"`
}
