package docker

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
)

func GetInnerRepoAddr(repo, operatorID, TaskName, localRegistry string) string {
	repository := repo
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", operatorID, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(localRegistry), repository, tag))
}

func Login(registry, username, password string) error {
	login := exec.Command("docker", "login", "-u", username, "-p", password, registry)
	output, err := login.CombinedOutput()
	if err != nil {
		return errors.Errorf("docker login failed, registry: %s, username: %s, err: %v", registry, username, string(output))
	}
	return nil
}
