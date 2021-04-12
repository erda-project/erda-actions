package build

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/java/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/dice"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/pack"
	"github.com/erda-project/erda-actions/pkg/render"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
)

const (
	// components 放置位置
	compPrefix    = "/opt/action/comp"
	cacheRootPath = "/opt/build_cache"
)

type JDKConfig struct {
	JavaHome  string
	SwitchCmd []string
}

var jdkSwitchCmdMap = map[string]*JDKConfig{
	"8": {
		JavaHome: "/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.272.b10-1.el7_9.x86_64",
		SwitchCmd: []string{
			"alternatives --set java /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.272.b10-1.el7_9.x86_64/jre/bin/java",
			"alternatives --set javac /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.272.b10-1.el7_9.x86_64/bin/javac",
		},
	},
	"11": {
		JavaHome: "/usr/lib/jvm/java-11-openjdk-11.0.6.10-1.el7_7.x86_64",
		SwitchCmd: []string{
			"alternatives --set java /usr/lib/jvm/java-11-openjdk-11.0.6.10-1.el7_7.x86_64/bin/java",
			"alternatives --set javac /usr/lib/jvm/java-11-openjdk-11.0.6.10-1.el7_7.x86_64/bin/javac",
		},
	},
}

func Execute() error {
	// 加载环境变量配置，配置来源: 1. 用户指定 2. pipeline指定
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	// 有缓存挂载目录,启用缓存
	if PathExists(cacheRootPath) {
		cacheKeyName := fmt.Sprintf("%s/%s/%s/%s", cfg.OrgName, cfg.ProjectName, cfg.AppName,
			strings.Replace(cfg.GittarBranch, "/", "-", -1))
		cacheStoragePath := cacheRootPath + "/" + cacheKeyName
		if !PathExists(cacheStoragePath) {
			err := os.MkdirAll(cacheStoragePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create dir %s err: %s", cacheStoragePath, err)
			}
		}
		cacheDirs := []string{"/root/.m2"}
		cacheMap := map[string]string{}
		for _, cacheDir := range cacheDirs {
			cacheFileName := base64.StdEncoding.EncodeToString([]byte(cacheDir)) + ".tar"
			cacheMap[cacheDir] = path.Join(cacheStoragePath, cacheFileName)
		}

		for cachePath, cacheTarFile := range cacheMap {
			// 之前运行存在缓存
			if PathExists(cacheTarFile) {
				os.MkdirAll(cachePath, os.ModePerm)
				fmt.Printf("start restore cache  %s\n", cacheTarFile)
				err := pack.UnTar(cacheTarFile, cachePath)
				if err != nil {
					return fmt.Errorf("failed to untar %s=>%s err: %s", cacheTarFile, cachePath, err)
				} else {
					fmt.Printf("restore cache  %s=>%s\n", cacheTarFile, cachePath)
				}
			} else {
				fmt.Printf("cacheFile:%s not exist\n", cacheTarFile)
			}
		}
		defer func() {
			for cacheDir, cacheTarFile := range cacheMap {
				if PathExists(cacheDir) {
					err := os.Chdir(cacheDir)
					fmt.Printf("pack cacheDir %s\n", cacheDir)
					if err != nil {
						fmt.Printf("failed to change cacheDir %s err: %s\n", cacheDir, err)
					}
					tmpFile := "/tmp/cache.tar"
					err = pack.Tar(tmpFile, ".")
					if err != nil {
						fmt.Printf("failed to tar %s=>%s err: %s\n", cacheDir, cacheTarFile, err)
					} else {
						fmt.Printf("success save cacheDir to %s", tmpFile)
					}
					mvCmd := exec.Command("mv", "-f", tmpFile, cacheTarFile)
					mvCmd.Stdout = os.Stdout
					mvCmd.Stdout = os.Stderr
					err = mvCmd.Run()
					if err != nil {
						fmt.Printf("failed to move %s=>%s err: %s\n", tmpFile, cacheTarFile, err)
					}
				} else {
					fmt.Printf("cacheDir %s not exist\n", cacheDir)
				}
			}
		}()
	}

	if cfg.ContainerType == "spring-boot" {
		// compatible with old spring-boot
		cfg.ContainerType = "openjdk"
	}
	jdkVersion := "8"
	if cfg.JDKVersion != nil {
		jdkVersion = fmt.Sprintf("%v", cfg.JDKVersion)
	}

	jdkConfig, ok := jdkSwitchCmdMap[jdkVersion]
	if !ok {
		return fmt.Errorf("not support java version %s", jdkVersion)
	}
	for _, switchCmd := range jdkConfig.SwitchCmd {
		err := runCommand(switchCmd)
		if err != nil {
			return err
		}
	}

	runCommand("export JAVA_HOME=" + jdkConfig.JavaHome)
	runCommand("java -version")

	fmt.Fprintln(os.Stdout, "successfully loaded action config")

	// 替换 maven settings & Dockerfile 占位符
	cfgMap := make(map[string]string)
	cfgMap["CENTRAL_REGISTRY"] = cfg.CentralRegistry
	cfgMap["NEXUS_URL"] = strutil.Concat("http://", strings.TrimPrefix(cfg.NexusAddr, "http://"))
	cfgMap["NEXUS_USERNAME"] = cfg.NexusUsername
	cfgMap["NEXUS_PASSWORD"] = cfg.NexusPassword
	if err := render.RenderTemplate(compPrefix, cfgMap); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "successfully replaced action placeholder")

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	// do build
	if err := build(cfg); err != nil {
		return err
	}

	// docker build & docker push 业务镜像
	if err := packAndPushAppImage(cfg); err != nil {
		return err
	}

	return nil
}

func build(cfg conf.Conf) error {
	// 切换工作目录
	if cfg.Context != "" {
		fmt.Fprintf(os.Stdout, "change workding directory to: %s\n", cfg.Context)
		if err := os.Chdir(cfg.Context); err != nil {
			return err
		}
	}
	fmt.Fprintf(os.Stdout, "build %s app\n", cfg.BuildType)
	switch cfg.BuildType {
	case "none":
		// nothing to do
	case "maven":
		mem := cfg.Memory
		if mem > 32 {
			mem = mem - 32
		}
		mvnOpts := fmt.Sprintf("-Xmx%dm", mem)
		if cfg.BuildCmd == "" {
			cfg.BuildCmd = "mvn clean package -e -B -U -Dmaven.test.skip"
		}
		mvnCmd := fmt.Sprintf("MAVEN_OPTS=%s %s -s %s %s",
			mvnOpts, cfg.BuildCmd, "/opt/action/comp/maven/settings.xml", cfg.Options)
		if err := simpleRun("/bin/bash", "-c", mvnCmd); err != nil {
			return err
		}
	case "gradle":
		mem := cfg.Memory
		if mem > 32 {
			mem = mem - 32
		}
		gradleOpts := fmt.Sprintf("-Xmx%dm", mem)
		// try to chmod
		_ = simpleRun("chmod", "+x", "./gradlew")
		if cfg.BuildCmd == "" {
			cfg.BuildCmd = "./gradlew build"
		}
		gradleCmd := fmt.Sprintf("GRADLE_OPTS=%s %s --init-script %s",
			gradleOpts, cfg.BuildCmd, "/opt/action/comp/gradle/init.gradle")
		if err := simpleRun("/bin/bash", "-c", gradleCmd); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown build_type: %s\n", cfg.BuildType)
	}
	fmt.Fprintln(os.Stdout, "successfully build")

	// url target 先下载
	if strings.Index(cfg.Target, "http://") == 0 || strings.Index(cfg.Target, "https://") == 0 {
		downloadTargetPath := "/tmp/target" + path.Ext(cfg.Target)
		err := dice.DownloadFile(cfg.Target, downloadTargetPath)
		if err != nil {
			return err
		}
		cfg.Target = downloadTargetPath
	}

	// 校验 target 是否存在
	if _, err := os.Stat(cfg.Target); err != nil {
		if os.IsNotExist(err) {
			// TODO: more friendly error message
		}
		return err
	}
	targetDir := filepath.Join(cfg.WorkDir, "target")
	if err := simpleRun("mkdir", "-p", targetDir); err != nil {
		return err
	}
	targetExt := filepath.Ext(cfg.Target)
	switch targetExt {
	case ".jar", ".war":
		// 拷贝 target 到 workdir
		if err := cp(cfg.Target, filepath.Join(targetDir, "app"+targetExt)); err != nil {
			return err
		}
	case ".tar":
		// gradle distribution
		distDir := filepath.Join(targetDir, "app")
		if err := simpleRun("mkdir", "-p", distDir); err != nil {
			return err
		}
		if err := simpleRun("tar", "-xvf", cfg.Target, "-C", distDir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown target file extension: %s\n", targetExt)
	}
	fmt.Fprintln(os.Stdout, "copy target success")
	if cfg.SwaggerPath != "" {
		if cfg.ServiceName == "" {
			return errors.New("need service_name param")
		}
		sjson, err := ioutil.ReadFile(cfg.SwaggerPath)
		if err != nil {
			return err
		}
		err = checkCompatibility(cfg, sjson)
		if err != nil {
			return err
		}
		// 拷贝 swagger 到 workdir
		if err := cp(cfg.SwaggerPath, filepath.Join(cfg.WorkDir, cfg.SwaggerPath)); err != nil {
			return err
		}
	}

	return nil
}

// docker build & docker push 业务镜像
func packAndPushAppImage(cfg conf.Conf) error {
	// 切换工作目录
	if err := os.Chdir(cfg.WorkDir); err != nil {
		return err
	}
	// copy assets
	if err := cp(compPrefix, "."); err != nil {
		return err
	}
	if cfg.Assets != "" {
		if err := mustDir(cfg.Assets); err != nil {
			return err
		}
		// copy assets
		if err := cp(cfg.Assets, "./assets"); err != nil {
			return err
		}
	} else {
		// no assets exist, but need create an empty dir
		if err := os.Mkdir("./assets", 0755); err != nil {
			return err
		}
	}
	// check container exist
	ct := fmt.Sprintf("%s/%s", compPrefix, cfg.ContainerType)
	if err := mustDir(ct); err != nil {
		fmt.Fprintf(os.Stdout, "container type: %s not exist", cfg.ContainerType)
		return err
	}
	dockerFilePath := fmt.Sprintf("%s/%s/Dockerfile", compPrefix, cfg.ContainerType)

	dockerCopyCmds := []string{}
	if len(cfg.CopyAssets) > 0 {
		rand := random.New()
		for _, asset := range cfg.CopyAssets {
			idx := strings.Index(asset, ":")
			if idx > 0 {
				source := asset[0:idx]
				if len(asset) == idx {
					// 模板文件路径为空不处理
					fmt.Fprintf(os.Stdout, "invalid asset: %s ", asset)
					continue
				}
				dest := asset[idx+1:]
				absSource := path.Join(cfg.Context, source)
				if strings.Index(source, "/") == 0 {
					//绝对路径
					absSource = source
				}
				sourceTarget := rand.String(6, random.Alphanumeric)
				if err := cp(absSource, sourceTarget); err != nil {
					return err
				}

				dockerCopyCmds = append(dockerCopyCmds, fmt.Sprintf("COPY %s %s", sourceTarget, dest))
			} else {
				if err := cp(path.Join(cfg.Context, asset), "target/"+asset); err != nil {
					return err
				}
			}
		}
	}
	if len(dockerCopyCmds) > 0 {
		dockerFileBytes, _ := ioutil.ReadFile(dockerFilePath)
		for _, cmd := range dockerCopyCmds {
			dockerFileBytes = append(dockerFileBytes, []byte("\n"+cmd)...)
		}
		ioutil.WriteFile(dockerFilePath, dockerFileBytes, os.ModePerm)
	}

	// jar包生成 & docker build 出业务镜像
	repo := getRepo(cfg)
	packCmd := exec.Command("docker", "build",
		"--build-arg", fmt.Sprintf("TARGET=%s", "target"),
		"--build-arg", fmt.Sprintf("MONITOR_AGENT=%s", cfg.MonitorAgent),
		"--build-arg", fmt.Sprintf("SPRING_PROFILES_ACTIVE=%s", cfg.Profile), // TODO: 非 spring 定制
		"--build-arg", fmt.Sprintf("DICE_VERSION=%s", cfg.DiceVersion),
		"--cpu-quota", strconv.FormatFloat(float64(cfg.CPU*100000), 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(cfg.Memory*apistructs.MB), 10),
		"-t", repo,
		"-f", fmt.Sprintf("%s/%s/Dockerfile", compPrefix, cfg.ContainerType), ".")
	if cfg.ContainerVersion != "" {
		packCmd.Args = append(packCmd.Args, "--build-arg", fmt.Sprintf("CONTAINER_VERSION=%s", cfg.ContainerVersion))
	}

	// 获取用户指定脚本命令
	scriptFile, err := os.Create("pre_start.sh")
	defer scriptFile.Close()
	if err != nil {
		return err
	}

	if cfg.PreStartScript != "" {
		fileInfo, err := os.Stat(cfg.PreStartScript)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return errors.New("is directory, not pre start script file")
		}

		input, err := ioutil.ReadFile(cfg.PreStartScript)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile("pre_start.sh", input, 0644)
		if err != nil {
			return err
		}
	}

	packCmd.Args = append(packCmd.Args, "--build-arg", fmt.Sprintf("SCRIPT_ARGS=%s", cfg.PreStartArgs))
	packCmd.Args = append(packCmd.Args, "--build-arg", fmt.Sprintf("WEB_PATH=%s", cfg.WebPath))

	fmt.Fprintf(os.Stdout, "packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully build app image: %s\n", repo)

	// docker push 业务镜像至集群 registry
	if err := simpleRun("docker", "push", repo); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully push app image: %s\n", repo)

	// upload metadata
	if err := storeMetaFile(&cfg, repo); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully upload metafile\n")
	if cfg.Service != "" {
		// TODO Deprecated: 使用 ${java:OUTPUT:image}
		// 写应用镜像信息至 pack-result, 供 release action 读取 & 填充dice.yml
		imageResult := make([]pack.ModuleImage, 0)
		imageResult = append(imageResult, pack.ModuleImage{ModuleName: cfg.Service, Image: repo})
		resultBytes, err := json.MarshalIndent(imageResult, "", "  ")
		if err != nil {
			return err
		}
		if err := filehelper.CreateFile(filepath.Join(cfg.WorkDir, "pack-result"), string(resultBytes), 0644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "successfully write image action: %s\n", repo)
	}

	if err := simpleRun("rm", "-rf", "comp", "target", "assets"); err != nil {
		fmt.Fprintf(os.Stdout, "warning, cleanup failed: %v", err)
	}

	return nil
}

// 生成业务镜像名称
func getRepo(cfg conf.Conf) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", cfg.TaskName, time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}

func storeMetaFile(cfg *conf.Conf, image string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "image",
				Value: image,
			},
		},
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(cfg.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func mustDir(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		// must a dir
		return fmt.Errorf("%s not a dir", path)
	}
	return nil
}

func cp(a, b string) error {
	fmt.Fprintf(os.Stdout, "copying %s to %s\n", a, b)
	return simpleRun("cp", "-r", a, b)
}

func simpleRun(name string, arg ...string) error {
	fmt.Fprintf(os.Stdout, "Run: %s, %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommand(cmd string) error {
	command := exec.Command("/bin/bash", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

func checkCompatibility(cfg conf.Conf, swagger []byte) error {
	var swaggerJson interface{}
	err := json.Unmarshal(swagger, &swaggerJson)
	if err != nil {
		return err
	}
	apiCheck := ApiCheck{
		OrgId:       strconv.Itoa(int(cfg.OrgID)),
		ProjectId:   strconv.Itoa(int(cfg.ProjectID)),
		AppId:       strconv.Itoa(int(cfg.AppID)),
		Workspace:   cfg.Workspace,
		ServiceName: cfg.ServiceName,
		ClusterName: cfg.ClusterName,
		RuntimeName: cfg.GittarBranch,
		Swagger:     swaggerJson,
	}

	body, err := json.Marshal(apiCheck)
	if err != nil {
		return err
	}
	headers := make(map[string]string)
	url := cfg.DiceOpenapiPrefix + "/api/gateway/check-compatibility"
	headers["Authorization"] = cfg.CiOpenapiToken
	rsp, _, err := Request("POST", url, body, 60, headers)
	if err != nil {
		return err
	}
	var result HttpResponse
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}
	if result.Success {
		return nil
	} else {
		return errors.New(result.Err.Msg)
	}

}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
