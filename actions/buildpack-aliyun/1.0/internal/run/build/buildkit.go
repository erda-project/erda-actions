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

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/build/buildcache"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda-actions/pkg/dockerfile"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
)


func BuildkitBuild() error {

	// before_build
	if err := beforeBuildForBuildkit(); err != nil {
		return errors.Wrap(err, "before build")
	}

	// build
	bplog.Println("开始制作编译镜像 ......")
	if err := dockerBuildForBuildkit(); err != nil {
		return errors.Wrap(err, "docker build Dockerfile")
	}

	// after_build
	if err := afterBuildForBuildkit(); err != nil {
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
func dockerBuildForBuildkit() error {

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

	memory := conf.PlatformEnvs().Memory

	mavenOpts := fmt.Sprintf("-Xmx%sm", strconv.FormatFloat(float64(memory-32), 'f', 0, 64))

	dockerBuildNetwork := os.Getenv("BP_DOCKER_BUILD_NETWORK")
	if dockerBuildNetwork == "" {
		dockerBuildNetwork = "host"
	}

	var nodeAdditionArgs []string
	if conf.Params().Language == types.LanguageNode {
		nodeAdditionArgs = []string{"--opt", "build-arg:" + fmt.Sprintf("NODE_OPTIONS=--max_old_space_size=%s",
			strconv.FormatFloat(float64(memory-32), 'f', 0, 64))}
	}
	//TODO: buildkitd addr from env
	buildCmdArgs := []string{
		"--addr",
		"tcp://buildkitd.default.svc.cluster.local:1234",
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--opt", "build-arg:PARENT_POM_DIR=/.cache_pom/parent_pom",
		"--opt", "build-arg:ALL_POM_DIR=/.cache_pom/all_pom",
		"--opt", "build-arg:" + fmt.Sprintf("FORCE_UPDATE_SNAPSHOT=%d", time.Now().Unix()),
		"--opt", "build-arg:MAVEN_OPTS=" + mavenOpts,
		"--opt", "build-arg:PACKAGE_LOCK_DIR=/.cache_packagejson",
		"--opt", "build-arg:DICE_VERSION=" + conf.PlatformEnvs().DiceVersion,
	}

	// HTTP_PROXY & HTTPS_PROXY
	if conf.Params().HttpProxy != "" {
		buildCmdArgs = append(buildCmdArgs, "--opt", "build-arg:HTTP_PROXY=" + conf.Params().HttpProxy)
	}
	if conf.Params().HttpsProxy != "" {
		buildCmdArgs = append(buildCmdArgs, "--opt", "build-arg:HTTPS_PROXY=" + conf.Params().HttpsProxy)
	}

	if len(nodeAdditionArgs) > 0 {
		buildCmdArgs = append(buildCmdArgs, nodeAdditionArgs...)
	}

	buildCmdArgs = append(buildCmdArgs,
		"--local", "context=" + conf.PlatformEnvs().WorkDir,
		"--local", "dockerfile=" + filepath.Join(conf.PlatformEnvs().WorkDir, "bp", "build"),
		"--output", "type=image,name=" + conf.EasyUse().DockerImageFromBuild + ",push=true,registry.insecure=true",
		"--import-cache", "type=registry,ref=" + conf.EasyUse().CalculatedCacheImage,
		"--export-cache", "type=registry,ref=" + conf.EasyUse().CalculatedCacheImage + ",push=true",
	)

	// build
	buildkitCmd := exec.Command("buildctl", buildCmdArgs...)
	bplog.Println(strutil.Join(buildkitCmd.Args, " ", false))
	bplog.Printf("build network: %s\n", dockerBuildNetwork)

	buildkitCmd.Dir = conf.PlatformEnvs().WorkDir
	buildkitCmd.Stdout = os.Stdout
	buildkitCmd.Stderr = os.Stderr
	if err := buildkitCmd.Run(); err != nil {
		return err
	}

	// 上报缓存镜像
	buildcache.ReportCacheImage("push")

	return nil
}

// execute dir is context
func beforeBuildForBuildkit() error {
	err := runPrepareScriptForBuildkit()
	if err != nil {
		return err
	}
	return nil
}

func runPrepareScriptForBuildkit() error {
	var script = []string{
		"#!/bin/sh",
		"set -eo pipefail",
		"w",
		"env | sort | grep -v USERNAME | grep -v PASSWORD || :",
		"free -h || :",
	}
	if conf.Params().BuildType == types.BuildTypeMaven || conf.Params().BuildType == types.BuildTypeMavenEdas {
		script = append(script, beforeBuildMavenForBuildkit()...)
	}
	if conf.Params().BuildType == types.BuildTypeNpm {
		script = append(script, beforeBuildNodeForBuilckit()...)
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
	return nil
}

func beforeBuildMavenForBuildkit() []string {
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

func beforeBuildNodeForBuilckit() []string {

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

func afterBuildForBuildkit() error {
	var buildResult = make([]ModuleArtifact, 0)
	var err error

	if conf.Params().Language == types.LanguageJava {
		buildResult, err = afterBuildJavaForBuildkit()
		if err != nil {
			return err
		}
	} else if conf.Params().BuildType == types.BuildTypeNpm {
		buildResult, err = afterBuildNodeForBuiltkit()
	}

	b, err := json.MarshalIndent(buildResult, "", "  ")
	if err != nil {
		return err
	}

	if err := filehelper.CreateFile(filepath.Join(conf.PlatformEnvs().WorkDir, "build-result"), string(b), 0644); err != nil {
		return err
	}
	return nil
}

// return: script, build-result, error
func afterBuildJavaForBuildkit() ([]ModuleArtifact, error) {
	artifactName := "app.jar"

	// build-result
	resName := filepath.Base(conf.PlatformEnvs().WorkDir)

	buildResult := make([]ModuleArtifact, 0)
	for _, m := range conf.Params().Modules {
		buildResult = append(buildResult, ModuleArtifact{m.Name, filepath.Join(resName, "app", m.Path, artifactName)})
	}

	return buildResult, nil
}

func afterBuildNodeForBuiltkit() ([]ModuleArtifact, error) {

	// build-result
	resName := filepath.Base(conf.PlatformEnvs().WorkDir)

	buildResult := make([]ModuleArtifact, 0)
	for _, m := range conf.Params().Modules {
		buildResult = append(buildResult, ModuleArtifact{m.Name, filepath.Join(resName, "app")})
	}

	return buildResult, nil
}

