package pkg

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	// env key of repo
	RepoToolsPath   = "REPO_TOOLS_PATH"
	RepoReleasePath = "REPO_RELEASE_PATH"
	RepoVersionPath = "REPO_VERSION_PATH"

	// env key of erda version
	DiceVersion = "DICE_VERSION"
	ErdaVersion = "ERDA_VERSION"

	// env key of erda public
	ErdaToPublic = "ERDA_TO_PUBLIC"

	// env key of git info
	GitAccount = "GIT_ACCOUNT"
	GitToken   = "GIT_TOKEN"

	// string of bool value
	True  = "true"
	False = "false"
)

// InitEnv set action env
func InitEnv(envs map[string]string) error {

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("set env %s=%s failed", k, v))
		}
	}

	return nil
}
