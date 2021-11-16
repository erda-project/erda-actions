package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
)

const (
	BUILDKIT_ENABLE = "BUILDKIT_ENABLE"
)

type AuthConfig struct {
	// key: registry address, value: auth info
	Auths map[string]RegistryAuthInfo `json:"auths"`
}

type RegistryAuthInfo struct {
	// auth info format: base64(username:password)
	Auth string `json:"auth"`
}

func GetInnerRepoAddr(repo, operatorID, TaskName, localRegistry string) string {
	repository := repo
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", operatorID, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(localRegistry), repository, tag))
}

func Login(registry, username, password string) error {
	if os.Getenv(BUILDKIT_ENABLE) == "true" {
		return GenerateAuthConfig(registry, username, password)
	}

	login := exec.Command("docker", "login", "-u", username, "-p", password, registry)
	output, err := login.CombinedOutput()
	if err != nil {
		return errors.Errorf("docker login failed, registry: %s, username: %s, err: %v", registry, username, string(output))
	}
	return nil
}

func GenerateAuthConfig(registry, username, password string) error {
	authConfigDir := os.Getenv("DOCKER_CONFIG")
	if authConfigDir == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get current user home dir failed: %v", err)
		}

		authConfigDir = fmt.Sprintf("%s/.docker", userHomeDir)
	}

	authConfigPath := fmt.Sprintf("%s/%s", authConfigDir, "config.json")

	// create docker config dir if it doesn't exist.
	authConfig, err := ioutil.ReadFile(authConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat docker config file, err: %v", err)
		}
		if err := os.MkdirAll(authConfigDir, 0755); err != nil {
			return fmt.Errorf("failed to create docker config dir, err: %v", err)
		}
	}

	ac := AuthConfig{
		Auths: map[string]RegistryAuthInfo{},
	}

	// unmarshal current docker config auth info if existed.
	if len(authConfig) != 0 {
		if err := json.Unmarshal(authConfig, &ac); err != nil {
			return fmt.Errorf("failed to unmarshal docker config, err: %v", err)
		}
	}

	// append new auth info
	ac.Auths[registry] = RegistryAuthInfo{
		Auth: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))),
	}

	// marshal new auth info
	configJson, err := json.Marshal(ac)
	if err != nil {
		return fmt.Errorf("marshal docker config json error: %v", err)
	}

	if err := ioutil.WriteFile(authConfigPath, configJson, 0644); err != nil {
		return fmt.Errorf("failed to write docker config json, path: %s, err: %v", authConfigDir, err)
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
		newError := errors.String()
		if strings.Contains(newError, "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"") {
			newError = strings.ReplaceAll(newError, "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"", "Docker registry 正在 gc ，请耐心等待 gc 完成")
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
		newError := errors.String()
		if strings.Contains(newError, "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"") {
			newError = strings.ReplaceAll(newError, "error parsing HTTP 405 response body: invalid character 'M' looking for beginning of value: \"Method not allowed\\n\"", "Docker registry 正在 gc ，请耐心等待 gc 完成")
		}
		return fmt.Errorf("%v", newError)
	}
	fmt.Printf("推送缓存镜像成功：\n")
	return nil
}
