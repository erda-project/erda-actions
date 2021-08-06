package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CopyRepoToTmp copy repo specified by repoVolumePath to /tmp
func CopyRepoToTmp(repoVolumePath string) error {

	// source repo exists validate
	exists, err := IsDirExists(repoVolumePath)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("stats repo path %s stat", repoVolumePath))
	}
	if !exists {
		return fmt.Errorf("source repo path %s does not exists", repoVolumePath)
	}

	// cp volume repo to tmp
	if _, err := ExecCmd(os.Stdout, os.Stderr, "", "cp",
		"-a", repoVolumePath, "/tmp/"); err != nil {
		return errors.WithMessage(err, "cp repo version to /tmp/")
	}

	return nil
}

// ReplaceFile replace file with content
func ReplaceFile(content, path string, perm os.FileMode) error {
	logrus.Infof("start to write file %s...", path)

	// build script in tools
	if err := ioutil.WriteFile(path, []byte(content), perm); err != nil {
		logrus.Info(err)
		return err
	}

	logrus.Infof("start to write file %s success!!", path)
	return nil
}

// RepoName get repo name from repo path
func RepoName(repoVolumePath string) string {
	_, name := path.Split(repoVolumePath)

	return name
}

// RepoTmpPath get repo path in /tmp
func RepoTmpPath(repoVolumePath string) string {
	return path.Join("/tmp", RepoName(repoVolumePath))
}

// GetPatchPath get erda release info path in /tmp
func GetPatchPath(erdaVersion string) string {
	return path.Join("/tmp", erdaVersion)
}
