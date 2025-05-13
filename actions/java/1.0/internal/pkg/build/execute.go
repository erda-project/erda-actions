package build

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/java/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/dice"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/render"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/metadata"
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
		JavaHome: "/usr/lib/jvm/java-1.8.0",
		SwitchCmd: []string{
			"alternatives --set java $(alternatives --list | grep java_sdk_1.8.0 | awk '{print $3}' | head -n 1)/jre/bin/java",
			"alternatives --set javac $(alternatives --list | grep java_sdk_1.8.0 | awk '{print $3}' | head -n 1)/bin/javac",
		},
	},
	"11": {
		JavaHome: "/usr/lib/jvm/java-11",
		SwitchCmd: []string{
			"alternatives --set java $(alternatives --list | grep java_sdk_11  | awk '{print $3}' | head -n 1)/bin/java",
			"alternatives --set javac $(alternatives --list | grep java_sdk_11  | awk '{print $3}' | head -n 1)/bin/javac",
		},
	},
	"17": {
		JavaHome: "/usr/lib/jvm/java-17",
		SwitchCmd: []string{
			"alternatives --set java $(alternatives --list | grep java_sdk_17  | awk '{print $3}' | head -n 1)/bin/java",
			"alternatives --set javac $(alternatives --list | grep java_sdk_17  | awk '{print $3}' | head -n 1)/bin/javac",
		},
	},
	"21": {
		JavaHome: "/usr/lib/jvm/java-21",
		SwitchCmd: []string{
			"alternatives --set java $(alternatives --list | grep java_sdk_21  | awk '{print $3}' | head -n 1)/bin/java",
			"alternatives --set javac $(alternatives --list | grep java_sdk_21  | awk '{print $3}' | head -n 1)/bin/javac",
		},
	},
}

func Execute() error {
	// 加载环境变量配置，配置来源: 1. 用户指定 2. pipeline指定
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	// 有缓存挂载目录,启用缓存
	if PathExists(cacheRootPath) {
		if err := handleCache(cfg); err != nil {
			return err
		}
	}

	// 切换至对应的 JDK 版本用于编译
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

	runCommand("echo export JAVA_HOME=" + jdkConfig.JavaHome + " >> /root/.bashrc")
	runCommand("echo export JAVA_HOME=" + jdkConfig.JavaHome + " >> /home/dice/.bashrc")
	runCommand("echo JAVA_HOME=$JAVA_HOME")
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

	// 提前执行 docker login
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
		sjson, err := os.ReadFile(cfg.SwaggerPath)
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

// storeMetaFile store meta data
func storeMetaFile(cfg *conf.Conf, image string) error {
	meta := apistructs.ActionCallback{
		Metadata: metadata.Metadata{
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
	cmd.Env = NewEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewEnv() []string {
	env := []string{
		"PATH=/opt/go/bin:/go/bin:/opt/nodejs/bin:/opt/maven/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	}

	jdkVersion := os.Getenv("ACTION_JDK_VERSION")
	if jdkVersion == "11" {
		env = append(env, "JAVA_HOME=/usr/lib/jvm/java-11")
	} else if jdkVersion == "17" {
		env = append(env, "JAVA_HOME=/usr/lib/jvm/java-17")
	} else if jdkVersion == "21" {
		env = append(env, "JAVA_HOME=/usr/lib/jvm/java-21")
	} else {
		env = append(env, "JAVA_HOME=/usr/lib/jvm/java-1.8.0")
	}

	return env
}

func runCommand(cmd string) error {
	command := exec.Command("/bin/bash", "-c", cmd)
	command.Env = NewEnv()
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
