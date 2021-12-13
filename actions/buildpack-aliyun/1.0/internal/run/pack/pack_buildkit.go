package pack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

func PackForBuildkit() ([]byte, error) {

	// before_pack

	// docker build
	bplog.Println("开始制作最终业务镜像 ......")
	packResult, err := dockerPackBuildForBuildkit()
	if err != nil {
		return nil, err
	}


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
func dockerPackBuildForBuildkit() ([]byte, error) {

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

	//----------------------
	newDockerfileContentLines := strings.Split(string(newDockerfileContent), "\n")
	FinallyDockerfileContentLines := append([]string{fmt.Sprintf("FROM %s as buildstage", conf.EasyUse().DockerImageFromBuild)}, newDockerfileContentLines...)
	for k, v := range FinallyDockerfileContentLines {
		if strings.HasPrefix(v, "ADD /app") {
			if conf.Params().Language == types.LanguageJava {
				FinallyDockerfileContentLines[k] = fmt.Sprintf("COPY --from=buildstage /code %s", strings.Fields(v)[2])
			} else if conf.Params().BuildType == types.BuildTypeNpm {
				FinallyDockerfileContentLines[k] = fmt.Sprintf("COPY --from=buildstage /app %s", strings.Fields(v)[2])
			}
		}
	}
	//-------------------------
	//if err = filehelper.CreateFile(exactDockerfilePath, string(newDockerfileContent), 0644); err != nil {
	//	return nil, err
	//}
	if err = filehelper.CreateFile(exactDockerfilePath, strings.Join(FinallyDockerfileContentLines, "\n"), 0644); err != nil {
		return nil, err
	}

	oneImage := conf.Params().Modules[0].Image.Name

	dockerBuildNetwork := os.Getenv("BP_DOCKER_BUILD_NETWORK")
	if dockerBuildNetwork == "" {
		dockerBuildNetwork = "host"
	}

	buildCmdArgs := []string{
		"--addr",
		"tcp://buildkitd.default.svc.cluster.local:1234",
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--opt", "build-arg:DICE_VERSION=" + conf.PlatformEnvs().DiceVersion,
		"--local",
		"context=" + exactWd,
		"--local",
		"dockerfile=" + filepath.Dir(exactDockerfilePath),
		"--output",
		"type=image,name=" + oneImage + ",push=true",
	}
	// HTTP_PROXY & HTTPS_PROXY
	if conf.Params().HttpProxy != "" {
		buildCmdArgs = append(buildCmdArgs, "--opt", "build-arg:HTTP_PROXY=" + conf.Params().HttpProxy)
	}
	if conf.Params().HttpsProxy != "" {
		buildCmdArgs = append(buildCmdArgs, "--opt", "build-arg:HTTPS_PROXY=" + conf.Params().HttpsProxy)
	}

	// build
	buildkitCmd := exec.Command("buildctl", buildCmdArgs...)

	bplog.Println(strutil.Join(buildkitCmd.Args, " ", false))
	bplog.Printf("docker build network: %s\n", dockerBuildNetwork)

	buildkitCmd.Dir = wd
	buildkitCmd.Stdout = os.Stdout
	buildkitCmd.Stderr = os.Stderr

	if err := buildkitCmd.Run(); err != nil {
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
		//dockerfileForARGPath := filepath.Join(wd, "Dockerfile.build."+m.Name)
		dockerfileForARGPath := filepath.Join(wd, "build", m.Name, "Dockerfile")
		if err := filehelper.CreateFile(dockerfileForARGPath, strings.Join(dockerfileForARG, "\n"), 0755); err != nil {
			return nil, err
		}
		tagPushScript = append(tagPushScript,
			fmt.Sprintf("buildctl --addr tcp://buildkitd.default.svc.cluster.local:1234 --tlscacert=/.buildkit/ca.pem --tlscert=/.buildkit/cert.pem --tlskey=/.buildkit/key.pem build --frontend dockerfile.v0 --local context=. --local dockerfile=%s --output type=image,name=%s,push=true", filepath.Dir(dockerfileForARGPath), m.Image.Name),
		)

		packResult = append(packResult, ModuleImage{m.Name, m.Image.Name})
	}

	tagPushScriptPath := filepath.Join(wd, "repack_push.sh")
	if err := filehelper.CreateFile(tagPushScriptPath, strings.Join(tagPushScript, "\n"), 075); err != nil {
		return nil, err
	}
	// push image TODO: FIX RETURN VALUE ERROR
	//if err = docker.PushByShell(tagPushScriptPath, wd); err != nil {
	//	return nil, err
	//}
	docker.PushByShell(tagPushScriptPath, wd)
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