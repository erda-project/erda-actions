package build

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/erda-project/erda-actions/pkg/render"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"

	"github.com/erda-project/erda-actions/actions/android/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

const (
	gradleConfigPath = "/root/.gradle/"
	scriptPath       = "/opt/action/command.sh"
)

type JDKConfig struct {
	JavaHome  string
	SwitchCmd []string
}

var jdkSwitchCmdMap = map[string]*JDKConfig{
	"8": {
		JavaHome: "/usr/lib/jvm/java-8-openjdk-amd64",
		SwitchCmd: []string{
			"update-alternatives --set java /usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java",
			"update-alternatives --set javac /usr/lib/jvm/java-8-openjdk-amd64/bin/javac",
		},
	},
	"11": {
		JavaHome: "/usr/lib/jvm/java-11-openjdk-amd64",
		SwitchCmd: []string{
			"update-alternatives --set java /usr/lib/jvm/java-11-openjdk-amd64/bin/java",
			"update-alternatives --set javac /usr/lib/jvm/java-11-openjdk-amd64/bin/javac",
		},
	},
}

var cfg conf.Conf

func Execute() error {

	envconf.MustLoad(&cfg)
	if err := os.Chdir(cfg.Context); err != nil {
		return err
	}
	if len(cfg.Target) == 0 && len(cfg.Targets) == 0 {
		return errors.New("no target specified")
	}
	// pipelineID is bigger than the max versionCode(2100000000)
	// see: https://developer.android.com/studio/publish/versioning
	//if err := createAndroidBuildCfg(); err != nil {
	//	return err
	//}
	// 替换 gradle 配置
	cfgMap := make(map[string]string)
	cfgMap["NEXUS_URL"] = strutil.Concat("http://", strings.TrimPrefix(cfg.NexusUrl, "http://"))
	cfgMap["NEXUS_USERNAME"] = cfg.NexusUsername
	cfgMap["NEXUS_PASSWORD"] = cfg.NexusPassword
	if err := render.RenderTemplate(gradleConfigPath, cfgMap); err != nil {
		return err
	}

	jdkVersion := "8"
	if cfg.JDKVersion != "" {
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

	fmt.Fprintln(os.Stdout, fmt.Sprintf("begin build target"))
	scriptContent := setupScript(cfg.Commands)
	err := ioutil.WriteFile(scriptPath, []byte(scriptContent), os.ModePerm)
	if err != nil {
		return err
	}
	if err := runCommand(scriptPath); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, fmt.Sprintf("build target success"))

	if len(cfg.Target) > 0 {
		if err := cpTarget(cfg.Target); err != nil {
			return err
		}
	}
	for _, target := range cfg.Targets {
		if err := cpTarget(target); err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stdout, "target files")
	runCommand("ls ", cfg.WorkDir)

	return nil
}

func runCommand(cmd ...string) error {
	c := strings.Join(cmd, " ")
	buildCmd := exec.Command("/bin/bash", "-c", c)

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}
	return nil
}

func cpTarget(target string) error {
	if _, err := os.Stat(target); err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("target file %s not exist", target))
		}
		return err
	}
	if err := cp(target, cfg.WorkDir); err != nil {
		return err
	}
	return nil
}

func cp(srcDir, destDir string) error {
	fmt.Fprintf(os.Stdout, "copying %s to %s\n", srcDir, destDir)
	cpCmd := exec.Command("cp", "-r", srcDir, destDir)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	return cpCmd.Run()
}

func createAndroidBuildCfg() error {
	if err := filehelper.CreateFile(cfg.Context+"/mobileBuild.cfg", "PIPELINE_ID="+cfg.PipelineID, 0644); err != nil {
		return err
	}
	return cp(cfg.Context+"/mobileBuild.cfg", cfg.WorkDir)
}
