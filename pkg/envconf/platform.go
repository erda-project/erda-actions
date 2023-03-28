package envconf

import (
	"fmt"
	"os"
)

type PlatformParams struct {
	ProjectID    uint64 `env:"DICE_PROJECT_ID"`
	ProjectName  string `env:"DICE_PROJECT_NAME"`
	AppID        uint64 `env:"DICE_APPLICATION_ID"`
	AppName      string `env:"DICE_APPLICATION_NAME"`
	Workspace    string `env:"DICE_WORKSPACE"`
	GittarRepo   string `env:"GITTAR_REPO"`
	GittarBranch string `env:"GITTAR_BRANCH"`
	GittarCommit string `env:"GITTAR_COMMIT"`
	OperatorID   string `env:"DICE_OPERATOR_ID"`
	PipelineID   int64  `env:"PIPELINE_ID"`

	OrgID  uint64 `env:"DICE_ORG_ID"`
	UserID string `env:"DICE_USER_ID"`

	DiceClusterName string `env:"DICE_CLUSTER_NAME" required:"true"`
	DiceArch        string `env:"DICE_ARCH" default:"amd64"`

	// metafile
	MetaFile string `env:"METAFILE"`

	LogID string `env:"TERMINUS_DEFINE_TAG"`

	// used to invoke openapi
	OpenAPIAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	OpenAPIToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
}

func NewPlatformParams() (PlatformParams, error) {
	platform := PlatformParams{}
	err := Load(&platform)
	if err != nil {
		return platform, err
	}
	return platform, nil
}

// GetTargetPlatforms return target platforms used for docker or buildkit build
// default platform is linux/{amd64/arm64}, if user want build multi-arch images
// it is good for user to set env PLATFORMS in application/settings/environ
func GetTargetPlatforms() string {
	arch := "amd64"
	if targetArch := os.Getenv("DICE_ARCH"); targetArch != "" {
		arch = targetArch
	}
	targetPlatforms := fmt.Sprintf("linux/%s", arch)
	if customPlatforms := os.Getenv("PLATFORMS"); customPlatforms != "" {
		targetPlatforms = customPlatforms
	}
	return targetPlatforms
}
