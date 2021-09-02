package build

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/build/buildcache"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda-actions/pkg/dockerfile"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
)

var (
	DownloadedDockerCachedImage = false
)

func Build() error {

	// before_build
	if err := beforeBuild(); err != nil {
		return errors.Wrap(err, "before build")
	}

	// build
	bplog.Println("开始制作编译镜像 ......")
	if err := dockerBuild(); err != nil {
		return errors.Wrap(err, "docker build Dockerfile")
	}

	// after_build
	if err := afterBuild(); err != nil {
		return errors.Wrap(err, "after build")
	}

	return nil
}

/*
context/
	-- repo/
	-- bp-backend/
		-- bp/
		-- code/ (context)
		-- before_build.sh
		-- .cache_pom/
		-- .cache_packagejson/
*/
func dockerBuild() error {

	// dockerfile just pack, no build
	if conf.Params().Language == types.LanguageDockerfile {
		return nil
	}

	dockerfilePath := filepath.Join(conf.PlatformEnvs().WorkDir, "bp", "build", "Dockerfile")

	dockerfileContent, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		return err
	}
	bpArgs := make(map[string]string)
	for k, v := range conf.Params().BpArgs {
		bpArgs[k] = v
	}
	if DownloadedDockerCachedImage {
		bpArgs["DEP_CACHE_IMAGE"] = conf.EasyUse().CalculatedCacheImage
	} else {
		bpArgs["DEP_CACHE_IMAGE"] = conf.EasyUse().DefaultCacheImage
	}
	var paths []string
	for _, v := range conf.Params().Modules {
		paths = append(paths, v.Path)
	}
	bpArgs["MODULES"] = "-am -pl " + strings.Join(paths, ",")
	newDockerfileContent := dockerfile.ReplaceOrInsertBuildArgToDockerfile(dockerfileContent, bpArgs)
	if err = filehelper.CreateFile(dockerfilePath, string(newDockerfileContent), 0644); err != nil {
		return err
	}

	cpu := conf.PlatformEnvs().CPU
	memory := conf.PlatformEnvs().Memory

	cpuShares := cpu * 1024
	cpuPeriod := 100000
	//cpuQuota := cpu * float64(cpuPeriod) * 2 // 0.5 -> 1, 1 -> 2
	cpuQuota := 150000

	mavenOpts := fmt.Sprintf("-Xmx%sm", strconv.FormatFloat(float64(memory-32), 'f', 0, 64))

	dockerBuildNetwork := os.Getenv("BP_DOCKER_BUILD_NETWORK")
	if dockerBuildNetwork == "" {
		dockerBuildNetwork = "host"
	}

	var nodeAdditionArgs []string
	if conf.Params().Language == types.LanguageNode {
		nodeAdditionArgs = []string{"--build-arg", fmt.Sprintf("NODE_OPTIONS=--max_old_space_size=%s",
			strconv.FormatFloat(float64(memory-32), 'f', 0, 64))}
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

		"--build-arg", "PARENT_POM_DIR=/.cache_pom/parent_pom",
		"--build-arg", "ALL_POM_DIR=/.cache_pom/all_pom",
		"--build-arg", fmt.Sprintf("FORCE_UPDATE_SNAPSHOT=%d", time.Now().Unix()),
		"--build-arg", "MAVEN_OPTS=" + mavenOpts,
		"--build-arg", "PACKAGE_LOCK_DIR=/.cache_packagejson",
		"--build-arg", "DICE_VERSION=" + conf.PlatformEnvs().DiceVersion,
	}

	// HTTP_PROXY & HTTPS_PROXY
	if conf.Params().HttpProxy != "" {
		dockerBuildCmdArgs = append(dockerBuildCmdArgs, "--build-arg", "HTTP_PROXY="+conf.Params().HttpProxy)
	}
	if conf.Params().HttpsProxy != "" {
		dockerBuildCmdArgs = append(dockerBuildCmdArgs, "--build-arg", "HTTPS_PROXY="+conf.Params().HttpsProxy)
	}

	if len(nodeAdditionArgs) > 0 {
		dockerBuildCmdArgs = append(dockerBuildCmdArgs, nodeAdditionArgs...)
	}
	dockerBuildCmdArgs = append(dockerBuildCmdArgs,
		"-t", conf.EasyUse().DockerImageFromBuild,
		"--cache-from", conf.EasyUse().CalculatedCacheImage,
		"-f", filepath.Join(conf.PlatformEnvs().WorkDir, "bp", "build", "Dockerfile"),
		conf.PlatformEnvs().WorkDir)

	// build
	dockerBuild := exec.Command("docker", dockerBuildCmdArgs...)

	bplog.Println(strutil.Join(dockerBuild.Args, " ", false))
	bplog.Printf("docker build network: %s\n", dockerBuildNetwork)

	dockerBuild.Dir = conf.PlatformEnvs().WorkDir
	dockerBuild.Stdout = os.Stdout
	dockerBuild.Stderr = os.Stderr
	if err := dockerBuild.Run(); err != nil {
		return err
	}

	// tag cache image
	dockerTag := exec.Command("docker", "tag", conf.EasyUse().DockerImageFromBuild, conf.EasyUse().CalculatedCacheImage)
	dockerTag.Stdout = os.Stdout
	dockerTag.Stderr = os.Stderr
	if err := dockerTag.Run(); err != nil {
		bplog.Printf("ReTag 缓存镜像失败，请忽略。镜像：%s，失败原因：%v\n", conf.EasyUse().DockerImageFromBuild, err)
		return nil
	}
	bplog.Printf("ReTag 缓存镜像成功：%s -> %s\n", conf.EasyUse().DockerImageFromBuild, conf.EasyUse().CalculatedCacheImage)
	// push cache image
	if err = docker.PushByCmd(conf.EasyUse().CalculatedCacheImage, conf.PlatformEnvs().WorkDir); err != nil {
		return err
	}
	// 上报缓存镜像
	buildcache.ReportCacheImage("push")

	return nil
}

// execute dir is context
func beforeBuild() error {
	err := runPrepareScript()
	if err != nil {
		return err
	}
	return nil
}

func runPrepareScript() error {
	var script = []string{
		"#!/bin/sh",
		"set -eo pipefail",
		"w",
		"env | sort | grep -v USERNAME | grep -v PASSWORD || :",
		"free -h || :",
	}
	if conf.Params().BuildType == types.BuildTypeMaven || conf.Params().BuildType == types.BuildTypeMavenEdas {
		script = append(script, beforeBuildMaven()...)
	}
	if conf.Params().BuildType == types.BuildTypeNpm {
		script = append(script, beforeBuildNode()...)
	}
	scriptPath := filepath.Join(conf.PlatformEnvs().WorkDir, "before_build.sh")
	err := filehelper.CreateFile(scriptPath, strings.Join(script, "\n"), 0755)
	if err != nil {
		return err
	}
	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Dir = filepath.Dir(conf.PlatformEnvs().WorkDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// docker pull cacheImage
	cacheCmd := exec.Command("docker", "pull", conf.EasyUse().CalculatedCacheImage)
	cacheCmd.Dir = filepath.Dir(conf.PlatformEnvs().WorkDir)
	cacheCmd.Stdout = os.Stdout
	devNull, _ := os.Open(os.DevNull)
	cacheCmd.Stderr = devNull
	bplog.Printf("开始下载编译缓存镜像: %s\n", conf.EasyUse().CalculatedCacheImage)
	// ignore error
	if err := cacheCmd.Run(); err != nil {
		DownloadedDockerCachedImage = false
		bplog.Printf("下载编译缓存镜像失败，请忽略。失败原因: %v\n", err)
	} else {
		DownloadedDockerCachedImage = true
		bplog.Println("下载编译缓存镜像成功!")
		buildcache.ReportCacheImage("pull")
	}

	return nil
}

func beforeBuildMaven() []string {
	pomDir := filepath.Join(conf.PlatformEnvs().WorkDir, ".cache_pom")
	err := filehelper.CheckExist(pomDir, true)
	if err != nil {
		bplog.Printf("判断 %s 目录是否存在失败，失败原因：%v\n", pomDir, err)
		_ = os.RemoveAll(pomDir)
	} else {
		return nil
	}
	return []string{
		"cd " + conf.Params().Context,

		"pom_dir=" + pomDir,
		"mkdir -p ${pom_dir}",

		// parent pom
		`echo ">>> copy parent pom file ......"`,
		"parent_pom_dir=${pom_dir}/parent_pom",
		"mkdir -p ${parent_pom_dir} && cp -f pom.xml ${parent_pom_dir}/pom.xml",

		// all pom
		`echo ">>> copy all pom files ......"`,
		`all_pom_dir=${pom_dir}/all_pom`,
		"mkdir -p ${all_pom_dir}",
		`find . -name 'pom.xml' -exec cp --parents {} ${all_pom_dir} \;`,

		// list pom
		`echo ">>> list pom.xml in ${pom_dir}:"`,
		`find ${pom_dir} -type f -name 'pom.xml' | cut -d '/' -f 6-`,
	}
}

func beforeBuildNode() []string {

	packageDir := filepath.Join(conf.PlatformEnvs().WorkDir, ".cache_packagejson")

	var script []string
	script = append(script,
		"#!/bin/sh",

		"cd "+conf.Params().Context,

		"p_cache_dir="+packageDir,
		"mkdir -p ${p_cache_dir}",

		`echo "check package-lock.json && package.json && .npmrc"`,

		`PACKAGE_JSON_FILE="package.json"`,
		`PACKAGE_LOCK_FILE="package-lock.json"`,
		`NPMRC_FILE=".npmrc"`,

		"if [ ! -f ${PACKAGE_JSON_FILE} ]",
		"then",
		`echo "No [${PACKAGE_JSON_FILE}] found! We need it."`,
		"exit -1",
		"fi",

		"if [ ! -f ${PACKAGE_LOCK_FILE} ]",
		"then",
		`echo "No [${PACKAGE_LOCK_FILE}] found! We need it. Please run [npm i] to generate this file and try again."`,
		"exit -1",
		"fi",
		"echo OK",

		`echo "copy package.json && package-lock.json && .npmrc into cache dir"`,

		`find . -maxdepth 4 -type f -name "package.json" '!' -path "*/node_modules/*" -exec cp --parents {} ${p_cache_dir} \;`,
		`find . -maxdepth 4 -type f -name "package-lock.json" '!' -path "*/node_modules/*" -exec cp --parents {} ${p_cache_dir} \;`,
		`find . -maxdepth 4 -type f -name ".npmrc" '!' -path "*/node_modules/*" -exec cp --parents {} ${p_cache_dir} \;`,
		`find . -maxdepth 4 -type f -name "webpack_dll.config.js" '!' -path "*/node_modules/*" -exec cp --parents {} ${p_cache_dir} \;`,

		`echo "list dependency cache files in: ${p_cache_dir}"`,
		`ls -lAR ${p_cache_dir}`,
	)

	return script
}

func afterBuild() error {

	var script string
	var buildResult = make([]ModuleArtifact, 0)
	var err error

	if conf.Params().Language == types.LanguageJava {
		script, buildResult, err = afterBuildJava()
		if err != nil {
			return err
		}
	} else if conf.Params().BuildType == types.BuildTypeNpm {
		script, buildResult, err = afterBuildNode()
	}

	scriptPath := filepath.Join(conf.PlatformEnvs().WorkDir, "after_build.sh")
	if err := filehelper.CreateFile(scriptPath, script, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(buildResult, "", "  ")
	if err != nil {
		return err
	}

	if err := filehelper.CreateFile(filepath.Join(conf.PlatformEnvs().WorkDir, "build-result"), string(b), 0644); err != nil {
		return err
	}

	bplog.Println("从编译镜像中获取编译产物 ......")

	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// return: script, build-result, error
func afterBuildJava() (string, []ModuleArtifact, error) {
	artifactName := "app.jar"

	var cpApps []string
	for _, m := range conf.Params().Modules {
		distination := filepath.Join(conf.PlatformEnvs().WorkDir, "app", m.Path)
		cpApps = append(cpApps,
			fmt.Sprintf("mkdir -p %s; docker cp ${temp_container}:/code/%s %s ", distination, filepath.Join(m.Path, artifactName), distination),
		)
	}
	// script
	var getArtifactScript []string
	getArtifactScript = append(getArtifactScript,
		"#!/bin/sh",
		`echo "get app from cache image ......"`,
		"temp_container=$(docker container create "+conf.EasyUse().DockerImageFromBuild+")",
		strings.Join(cpApps, "\n"),
		// list compiled files
		`cd `+filepath.Join(conf.PlatformEnvs().WorkDir, "app")+` && find . -type f`,
		`docker container rm -f ${temp_container} >/dev/null 2>&1 || true`,
		`echo "done!"`,
	)

	// build-result
	resName := filepath.Base(conf.PlatformEnvs().WorkDir)

	buildResult := make([]ModuleArtifact, 0)
	for _, m := range conf.Params().Modules {
		buildResult = append(buildResult, ModuleArtifact{m.Name, filepath.Join(resName, "app", m.Path, artifactName)})
	}

	return strings.Join(getArtifactScript, "\n"), buildResult, nil
}

func afterBuildNode() (string, []ModuleArtifact, error) {

	// script
	var getArtifactScript []string
	getArtifactScript = append(getArtifactScript,
		"#!/bin/sh",
		`echo "get app from cache image ......"`,
		"temp_container=$(docker container create "+conf.EasyUse().DockerImageFromBuild+")",
		"docker cp ${temp_container}:/app/. "+filepath.Join(conf.PlatformEnvs().WorkDir, "app"),
		// list compiled files
		`cd `+filepath.Join(conf.PlatformEnvs().WorkDir, "app")+` && find . -type f '!' -path "*/node_modules/*"`,
		`docker container rm -f ${temp_container} >/dev/null 2>&1 || true`,
		`echo "done!"`,
	)

	// build-result
	resName := filepath.Base(conf.PlatformEnvs().WorkDir)

	buildResult := make([]ModuleArtifact, 0)
	for _, m := range conf.Params().Modules {
		buildResult = append(buildResult, ModuleArtifact{m.Name, filepath.Join(resName, "app")})
	}

	return strings.Join(getArtifactScript, "\n"), buildResult, nil
}

type ModuleArtifact struct {
	ModuleName   string `json:"module_name"`
	ArtifactPath string `json:"artifact_path"`
}
