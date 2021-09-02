package docker

import (
	"bytes"
	"fmt"
	"os"
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

func PushByCmd(imageName string, workDir string) error {
	var errors bytes.Buffer
	dockerPush := exec.Command("docker", "push", imageName)
	if workDir != "" {
		dockerPush.Dir = workDir
	}
	dockerPush.Stdout = os.Stdout
	dockerPush.Stderr = &errors
	if err := dockerPush.Run(); err != nil {
		fmt.Printf("推送缓存镜像失败，请忽略。镜像：%s，失败原因：%v\n\n", imageName, err)
	}

	// error 信息大于 0
	if errors.Len() > 0 {
		var newError string
		if strings.Contains(errors.String(), "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"") {
			newError = strings.ReplaceAll(errors.String(), "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"", "Docker registry 正在 gc ，请耐心等待 gc 完成")
		}
		return fmt.Errorf("%v", newError)
	}

	fmt.Printf("推送缓存镜像成功：%s\n", imageName)
	return nil
}

func PushByShell(pushScriptPath string, workDir string) error {
	var errors bytes.Buffer
	dockerPush := exec.Command("/bin/sh", pushScriptPath)
	if workDir != "" {
		dockerPush.Dir = workDir
	}
	dockerPush.Stdout = os.Stdout
	dockerPush.Stderr = &errors
	if err := dockerPush.Run(); err != nil {
		fmt.Printf("推送缓存镜像失败，请忽略。失败原因：%v\n\n", err)
	}

	// error 信息大于 0
	if errors.Len() > 0 {
		var newError string
		if strings.Contains(errors.String(), "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"") {
			newError = strings.ReplaceAll(errors.String(), "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"", "Docker registry 正在 gc ，请耐心等待 gc 完成")
		}
		return fmt.Errorf("%v", newError)
	}
	fmt.Printf("推送缓存镜像成功：\n")
	return nil
}
