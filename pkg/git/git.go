package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	// git remote name, instead of default origin
	GittarRemote = "gittar"
	TmpDir       = "/tmp/analyzer"
)

var logger = log.New(os.Stdout, "[Action] ", 0)

func FetchRepo(gitUrl, branch string, dest string) (string, error) {

	if strings.HasPrefix(gitUrl, "file://") {
		gitUrl = strings.TrimSuffix(gitUrl, ".git")
	}
	logger.Printf("Git Repo: %s, Branch: %s\n", gitUrl, branch)
	var projectDir string
	if dest == "" {
		projectDir = filepath.Join(TmpDir, gitUrl, branch)
	} else {
		projectDir = filepath.Clean(dest)
	}

	var gitCmd *exec.Cmd
	var errFormat string

	// git pull || git clone
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		logger.Printf("directory %s not exist, will auto create later\n", projectDir)
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return "", errors.Wrapf(err, "create directory error: %s", projectDir)
		}

		// git clone
		gitCmd = exec.Command(
			"git", "clone", "--recursive",
			"--single-branch", "--depth=1", "--branch", branch,
			"--origin", GittarRemote,
			"--",
			gitUrl, projectDir,
		)
		errFormat = "git clone error: %s"

	} else {
		gitCmd = exec.Command(
			"/bin/sh", "-c",
			fmt.Sprintf("git fetch --recurse-submodules %s && git reset --hard %s/%s",
				GittarRemote, GittarRemote, branch),
		)
		gitCmd.Dir = projectDir
		errFormat = "git pull error: %s"
	}

	if content, err := gitCmd.CombinedOutput(); err != nil {
		if err2 := os.RemoveAll(projectDir); err2 != nil {
			errFormat = fmt.Sprintf("%s; and fail to remove tmp directory: %s, cause: %v",
				errFormat, projectDir, err2)
		}
		return "", errors.Wrapf(err, errFormat, string(content))
	}

	return projectDir, nil
}
