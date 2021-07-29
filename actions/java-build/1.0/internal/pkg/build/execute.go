package build

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/erda-project/erda-actions/actions/java-build/1.0/internal/pkg/conf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/render"
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

	runCommand(fmt.Sprintf(" mkdir $MAVEN_CONFIG "))
	//runCommand(fmt.Sprintf(" cp %s %s ", "/opt/action/comp/maven/settings.xml", "$MAVEN_CONFIG/settings.xml"))
	runCommand(fmt.Sprintf(" cp %s %s ", "/opt/action/comp/maven/settings.xml", "$MAVEN_HOME/conf/settings.xml"))

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

	//做一些agent的工作，将dockerfile中下载和拷贝到agent.jar文件拷贝到buildPath目录下
	runCommand(fmt.Sprintf(" mkdir -p %s", fmt.Sprintf("%s/%s/%s", cfg.WorkDir, pwdName, "spot-agent")))
	runCommand(fmt.Sprintf(" cp -rv %s %s ", "/opt/action/comp/spot-agent/${DICE_VERSION}/spot-agent/.", fmt.Sprintf("%s/%s/%s", cfg.WorkDir, pwdName, "spot-agent")))
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

	result, _ := ioutil.ReadAll(out) // 读取输出结果
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
