package build

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/java-build/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/jdk"
	"github.com/erda-project/erda-actions/pkg/render"
	"github.com/erda-project/erda-actions/pkg/version"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/strutil"
)

const (
	// components 放置位置
	compPrefix = "/opt/action/comp"
)

type JDKConfig struct {
	JavaHome  string
	SwitchCmd []string
}

func Execute() error {
	// 加载环境变量配置，配置来源: 1. 用户指定 2. pipeline指定
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	// docker login
	if cfg.LocalRegistryUserName != "" {
		if err := docker.Login(cfg.LocalRegistry, cfg.LocalRegistryUserName, cfg.LocalRegistryPassword); err != nil {
			return err
		}
	}

	javaSwitcher := jdk.NewUpdateAlternativesSwitcher()
	jdkVersionStr := strconv.Itoa(cfg.JDKVersion)
	if !javaSwitcher.IsVersionSupported(jdkVersionStr) {
		return fmt.Errorf("unsupported Java version %d", cfg.JDKVersion)
	}

	if err := javaSwitcher.SwitchToVersion(jdkVersionStr); err != nil {
		return fmt.Errorf("failed to switch Java version: %v", err)
	}

	actualJavaHome, err := javaSwitcher.GetCurrentJavaHome()
	if err != nil {
		return fmt.Errorf("failed to get JAVA_HOME: %v", err)
	}

	runCommand("echo export JAVA_HOME=" + actualJavaHome + " >> /root/.bashrc")
	runCommand("echo export JAVA_HOME=" + actualJavaHome + " >> /home/dice/.bashrc")
	runCommand("export JAVA_HOME=" + actualJavaHome)
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

	// do build
	if err := build(cfg); err != nil {
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

	//runCommand(fmt.Sprintf(" mkdir $MAVEN_CONFIG "))
	//runCommand(fmt.Sprintf(" cp %s %s ", "/opt/action/comp/maven/settings.xml", "$MAVEN_CONFIG/settings.xml"))
	runCommand("mkdir -p /root/.m2")
	runCommand(fmt.Sprintf(" cp -f %s %s ", "/opt/action/comp/maven/settings.xml", "/root/.m2/settings.xml"))

	_ = simpleRun("chmod", "+x", "./gradlew")
	runCommand("mkdir /root/.gradle")
	runCommand(fmt.Sprintf(" cp %s %s ", "/opt/action/comp/gradle/init.gradle", "/root/.gradle/init.gradle"))

	for _, v := range cfg.BuildCmd {
		if err := simpleRun("/bin/bash", "-c", v); err != nil {
			return err
		}
	}

	if cfg.Context == "" {
		return errors.Errorf("Need to specify workdir")
	}

	//获取当前工作目录的目录名称
	pwdName := runCmdBackResult("basename `pwd`")
	pwdName = strings.Replace(pwdName, "\n", "", -1)
	if err := os.Chdir(cfg.WorkDir); err != nil {
		return err
	}
	//将工作目录中的文件都拷贝到当前action的目录下
	runCommand(fmt.Sprintf(" rm -rf %s", pwdName))
	runCommand(fmt.Sprintf(" mkdir %s", pwdName))
	runCommand(fmt.Sprintf(" cp -r %s ./", cfg.Context))

	//创建输出OUTPUT文件
	if !filepath.IsAbs(cfg.MetaFile) {
		return errors.Errorf("not an absolute path: %s", cfg.MetaFile)
	}
	err := os.MkdirAll(filepath.Dir(cfg.MetaFile), 0755)
	if err != nil {
		return errors.Wrap(err, "make parent dir error")
	}
	//输出buildPath到OUTPUT中
	if err := runCommand(fmt.Sprintf("echo 'buildPath=%s' >> %s ", fmt.Sprintf("%s/%s", cfg.WorkDir, pwdName), cfg.MetaFile)); err != nil {
		logrus.Errorf(" write buildPath error: ", err)
		return err
	}

	erdaVersion := cfg.DiceVersion
	if !version.IsHistoryVersion(erdaVersion) {
		erdaVersion = "latest"
	}

	// java agent version
	agentFileName := "spot-agent.tar.gz"
	if cfg.JDKVersion >= 17 {
		agentFileName = "spot-agent-jdk17.tar.gz"
	}

	//做一些agent的工作，将dockerfile中下载和拷贝到agent.jar文件拷贝到buildPath目录下
	runCommand(fmt.Sprintf("mkdir -p %s", fmt.Sprintf("%s/%s", cfg.WorkDir, pwdName)))
	runCommand(fmt.Sprintf("tar -xzf %s -C %s", fmt.Sprintf("/opt/action/comp/spot-agent/%s/%s", erdaVersion, agentFileName), fmt.Sprintf("%s/%s", cfg.WorkDir, pwdName)))
	runCommand(fmt.Sprintf("echo 'JAVA_OPTS=%s' >> %s ", "-javaagent:/spot-agent/spot-agent.jar", cfg.MetaFile))

	return nil
}

func runCmdBackResult(cmd string) string {
	command := exec.Command("/bin/bash", "-c", cmd)
	out, _ := command.StdoutPipe()
	defer func() {
		if out != nil {
			out.Close()
		}
	}()

	if err := command.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	result, _ := io.ReadAll(out)
	return string(result)
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
