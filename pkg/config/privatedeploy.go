package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type PrivateDeployConfig struct {
	// -- buildpack
	// common
	BpDockerBaseRegistry     string `env:"BP_DOCKER_BASE_REGISTRY"`
	BpDockerCacheRegistry    string `env:"BP_DOCKER_CACHE_REGISTRY"`
	BpDockerArtifactRegistry string `env:"BP_DOCKER_ARTIFACT_REGISTRY"`

	// java
	BpNexusURL      string `env:"BP_NEXUS_URL"`
	BpNexusUsername string `env:"BP_NEXUS_USERNAME"`
	BpNexusPassword string `env:"BP_NEXUS_PASSWORD"`

	// nodejs / spa
}

var privateDeployCfg PrivateDeployConfig
var privateDeployCfgMap = map[string]string{}
var hasParsed = false

func GetPrivateDeployConfig() PrivateDeployConfig {
	if !hasParsed {
		parsePrivateDeployConfig()
	}
	return privateDeployCfg
}

func GetPrivateDeployConfigMap() map[string]string {
	if !hasParsed {
		parsePrivateDeployConfig()
	}
	return privateDeployCfgMap
}

func parsePrivateDeployConfig() error {
	if hasParsed {
		return nil
	}
	if err := env.Parse(&privateDeployCfg); err != nil {
		return errors.Wrap(err, "parse envs of PRIVATE_DEPLOY")
	}
	privateDeployCfgMap["BP_DOCKER_BASE_REGISTRY"] = privateDeployCfg.BpDockerBaseRegistry
	privateDeployCfgMap["BP_DOCKER_CACHE_REGISTRY"] = privateDeployCfg.BpDockerCacheRegistry
	privateDeployCfgMap["BP_DOCKER_ARTIFACT_REGISTRY"] = privateDeployCfg.BpDockerArtifactRegistry
	privateDeployCfgMap["BP_NEXUS_URL"] = privateDeployCfg.BpNexusURL
	privateDeployCfgMap["BP_NEXUS_USERNAME"] = privateDeployCfg.BpNexusUsername
	privateDeployCfgMap["BP_NEXUS_PASSWORD"] = privateDeployCfg.BpNexusPassword

	privateDeployCfgMap["BP_DOCKER_REGISTRY"] = privateDeployCfg.BpDockerBaseRegistry //兼容2.12的buildpack配置
	hasParsed = true
	return nil
}
