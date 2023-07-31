package conf

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda/pkg/envconf"
)

// platformEnv 平台注入的环境变量
type platformEnv struct {
	GittarBranch         string `env:"GITTAR_BRANCH" required:"true"`
	ClusterName          string `env:"DICE_CLUSTER_NAME" required:"true"`
	OpenAPIToken         string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DefaultDepCacheImage string `env:"DEFAULT_DEP_CACHE_IMAGE"`
	ExportCacheType      string `env:"EXPORT_CACHE_TYPE" default:"registry"`
	PipelineID           uint64 `env:"PIPELINE_ID" required:"true"`
	ProjectAppAbbr       string `env:"DICE_PROJECT_APPLICATION"`
	DiceWorkspace        string `env:"DICE_WORKSPACE" required:"true"`
	ContextDir           string `env:"CONTEXTDIR" required:"true"`
	WorkDir              string `env:"WORKDIR" required:"true"`
	MetaFile             string `env:"METAFILE" required:"true"`
	DiceVersion          string `env:"DICE_VERSION" required:"true"`

	// resources
	CPU    float64 `env:"PIPELINE_LIMITED_CPU" required:"true"` // 核数
	Memory float64 `env:"PIPELINE_LIMITED_MEM" required:"true"` // 单位: MB

	// nexus
	NexusAddr     string `env:"NEXUS_ADDR"`
	NexusUserName string `env:"NEXUS_USERNAME"`
	NexusPassword string `env:"NEXUS_PASSWORD"`

	// docker registry
	DockerRegistry string `env:"DOCKER_REGISTRY"`
	// 构建产物应该推送至的 docker registry
	BpDockerArtifactRegistry         string `env:"BP_DOCKER_ARTIFACT_REGISTRY"`
	BpDockerArtifactRegistryUserName string `env:"BP_DOCKER_ARTIFACT_REGISTRY_USERNAME"`
	BpDockerArtifactRegistryPassword string `env:"BP_DOCKER_ARTIFACT_REGISTRY_PASSWORD"`
	// 缓存镜像应该推送至的 docker registry
	BpDockerCacheRegistry string `env:"BP_DOCKER_CACHE_REGISTRY"`
}

func initPlatformEnvs() (*platformEnv, error) {
	var env platformEnv
	if err := envconf.Load(&env); err != nil {
		return nil, errors.Errorf("failed to parse platform envs: %v", err)
	}

	if env.BpDockerArtifactRegistryUserName != "" {
		err := docker.Login(env.BpDockerArtifactRegistry, env.BpDockerArtifactRegistryUserName, env.BpDockerArtifactRegistryPassword)
		if err != nil {
			return nil, err
		}
	}
	switch env.ExportCacheType {
	case "inline", "registry", "local", "gha":
	default:
		return nil, errors.Errorf("invalid export cache type: %s", env.ExportCacheType)
	}

	return &env, nil
}
