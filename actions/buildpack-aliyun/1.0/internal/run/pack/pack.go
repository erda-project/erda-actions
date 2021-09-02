package pack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/dockerfile"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
	"github.com/erda-project/erda/pkg/template"
)

func Pack() ([]byte, error) {

	// before_pack

	// docker build
	bplog.Println("开始制作最终业务镜像 ......")
	packResult, err := dockerPackBuild()
	if err != nil {
		return nil, err
	}

	// after_pack

	return packResult, nil
}

/*
context/
	-- repo/
	-- bp-backend/
		-- bp/
		-- code/
		-- before_build.sh
		-- .cache_pom/
		-- .cache_packagejson/
		-- app/
*/
func dockerPackBuild() ([]byte, error) {

	wd := conf.PlatformEnvs().WorkDir

	// bp at wd/bp
	// copy code context to wd/code
	err := filehelper.CheckExist(filepath.Join(wd, "code"), true)
	if err != nil {
		if err = os.RemoveAll(filepath.Join(wd, "code")); err != nil {
			return nil, err
		}
		if err = copy.Copy(conf.Params().Context, filepath.Join(wd, "code")); err != nil {
			return nil, err
		}
	}

	var exactDockerfilePath = filepath.Join(wd, "bp", "pack", "Dockerfile")
	var startFilePath = filepath.Join(wd, "bp", "pack", "start.sh")
	startFileContent, err := ioutil.ReadFile(startFilePath)
	if err == nil {
		newStartFileContent := template.Render(string(startFileContent),
			map[string]string{"JAVA_OPTS": conf.Params().JavaOpts},
		)
		logrus.Infof("new startfile :%s", newStartFileContent)
		err := ioutil.WriteFile(startFilePath, []byte(newStartFileContent), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	var exactWd = wd
	// dockerfile use user's Dockerfile
	if conf.Params().Language == types.LanguageDockerfile {
		exactDockerfilePath = filepath.Join(wd, "code", "Dockerfile")
		exactWd = filepath.Join(wd, "code")
	}

	dockerfileContent, err := ioutil.ReadFile(exactDockerfilePath)
	if err != nil {
		return nil, err
	}
	newDockerfileContent := dockerfile.ReplaceOrInsertBuildArgToDockerfile(dockerfileContent, conf.Params().BpArgs)
	if err = filehelper.CreateFile(exactDockerfilePath, string(newDockerfileContent), 0644); err != nil {
		return nil, err
	}

	cpu := conf.PlatformEnvs().CPU
	memory := conf.PlatformEnvs().Memory

	cpuShares := cpu * 1024
	cpuPeriod := 100000
	//cpuQuota := cpu * float64(cpuPeriod) * 2 // 0.5 -> 1, 1 -> 2
	cpuQuota := 150000

	oneImage := conf.Params().Modules[0].Image.Name

	dockerBuildNetwork := os.Getenv("BP_DOCKER_BUILD_NETWORK")
	if dockerBuildNetwork == "" {
		dockerBuildNetwork = "host"
	}

	dockerBuildCmdArgs := []string{
		"build",
		// float
		"--memory", strconv.FormatFloat(float64(memory*1000000), 'f', 0, 64),
		// int strconv.ParseInt
		"--cpu-shares", strconv.FormatFloat(float64(cpuShares), 'f', 0, 64),
		// int
		"--cpu-period", strconv.FormatFloat(float64(cpuPeriod), 'f', 0, 64),
		// int
		"--cpu-quota", strconv.FormatFloat(float64(cpuQuota), 'f', 0, 64),
		"--network", dockerBuildNetwork,

		"--pull",

		"-t", oneImage,
		"-f", exactDockerfilePath,

		"--build-arg", "DICE_VERSION=" + conf.PlatformEnvs().DiceVersion,

		exactWd,
	}

	// HTTP_PROXY & HTTPS_PROXY
	if conf.Params().HttpProxy != "" {
		dockerBuildCmdArgs = append(dockerBuildCmdArgs, "--build-arg", "HTTP_PROXY="+conf.Params().HttpProxy)
	}
	if conf.Params().HttpsProxy != "" {
		dockerBuildCmdArgs = append(dockerBuildCmdArgs, "--build-arg", "HTTPS_PROXY="+conf.Params().HttpsProxy)
	}

	// build
	dockerBuild := exec.Command("docker", dockerBuildCmdArgs...)

	bplog.Println(strutil.Join(dockerBuild.Args, " ", false))
	bplog.Printf("docker build network: %s\n", dockerBuildNetwork)

	dockerBuild.Dir = wd
	dockerBuild.Stdout = os.Stdout
	dockerBuild.Stderr = os.Stderr
	if err := dockerBuild.Run(); err != nil {
		return nil, err
	}

	// 0. 给 image 打上 APP_DIR env
	// 1. multi module docker tag and push
	// 2. 写 pack-result 文件

	packResult := make([]ModuleImage, 0)
	var tagPushScript = []string{"#!/bin/sh"}
	for _, m := range conf.Params().Modules {
		dockerfileForARG := []string{
			fmt.Sprintf("FROM %s AS base", oneImage),
			fmt.Sprintf("ENV APP_DIR=%s", m.Path),
			fmt.Sprintf("FROM base"),
		}
		dockerfileForARGPath := filepath.Join(wd, "Dockerfile.build."+m.Name)
		if err := filehelper.CreateFile(dockerfileForARGPath, strings.Join(dockerfileForARG, "\n"), 0755); err != nil {
			return nil, err
		}
		tagPushScript = append(tagPushScript,
			fmt.Sprintf("docker build -t %s -f %s .", m.Image.Name, dockerfileForARGPath),
			fmt.Sprintf("docker push %s", m.Image.Name),
		)

		packResult = append(packResult, ModuleImage{m.Name, m.Image.Name})
	}

	tagPushScriptPath := filepath.Join(wd, "repack_push.sh")
	if err := filehelper.CreateFile(tagPushScriptPath, strings.Join(tagPushScript, "\n"), 075); err != nil {
		return nil, err
	}

	// push image
	if err = docker.PushByShell(tagPushScriptPath, wd); err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(packResult, "", "  ")
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(wd)
	if err != nil {
		return nil, err
	}
	if err := filehelper.CreateFile(filepath.Join(wd, "pack-result"), string(b), 0644); err != nil {
		return nil, err
	}

	return b, err
}

// afterPack will do something
func afterPack() error {
	return nil
}

type ModuleImage struct {
	ModuleName string `json:"module_name"`
	Image      string `json:"image"`
}
