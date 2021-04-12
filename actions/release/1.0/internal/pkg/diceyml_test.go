package pkg

import (
	"fmt"
	"testing"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/release/1.0/internal/diceyml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var envyml = `
addons:
  terminus-elasticsearch:
    plan: "terminus-elasticsearch:professional"
    options:
      version: "9.0"
  xxx:
    plan: "mysql:basic"
    options:
      version: "8.0"
`

var yml = `version: 2.0

version: 2
envs:
  TERMINUS_APP_NAME: "TEST-global"
  TEST_PARAM: "param_value"
services:
  web:
    ports:
      - 8080
      - port: 20880
      - port: 1234
        protocol: "UDP"
      - port: 4321
        protocol: "HTTP"
      - port: 53
        protocol: "DNS"
        l4_protocol: "UDP"
        default: true
    health_check:
      exec:
        cmd: "echo 1"
    deployments:
      replicas: ${replicas}
    resources:
      cpu: ${cpu:0.1}
      mem: 512
      disk: 0
    expose:
      - 20880
    volumes:
      - storage: "nfs"
        path: "/data/file/resource"
    endpoints:
      - domain: aaa.com
    traffic_security:
      mode: https
addons:
  terminus-elasticsearch:
    plan: "terminus-elasticsearch:professional"
    options:
      version: "6.8.9"
  xxx:
    plan: "mysql:basic"
    options:
      version: "5.7.23"
values:
  test:
    replicas: 1
    cpu: 0.5
  production:
    replicas: 2
    cpu: 1

`

func TestDiceYml(t *testing.T) {
	d, err := diceyml.New([]byte(yml))
	assert.Nil(t, err)
	envd, err := diceyml.New([]byte(envyml))
	assert.Nil(t, err)
	err = d.Compose("production", envd)
	err = serviceMeshAddonAdjust(d)
	assert.Nil(t, err)
	err = apiGatewayAddonAdjust(d)
	assert.Nil(t, err)
	cfg := &conf.Conf{
		Services: map[string]conf.Service{
			"web": conf.Service{
				Cmd: "run",
			},
		},
	}
	insertCommands(d, cfg)
	err = insertAddons(d, "PROD", "mysql", map[string]string{"init_sql": "hhh"})
	assert.Nil(t, err)
	d.SetEnv("zty", "cc")
	b, _ := yaml.Marshal(d.Obj())
	fmt.Println(string(b))
}
