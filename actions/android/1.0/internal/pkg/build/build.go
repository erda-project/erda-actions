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

var cfg conf.Conf

func Execute() error {

	envconf.MustLoad(&cfg)
	if err := os.Chdir(cfg.Context); err != nil {
		return err
	}
	if err := createAndroidBuildCfg(); err != nil {
		return err
	}
	// 替换 gradle 配置
	cfgMap := make(map[string]string)
	cfgMap["NEXUS_URL"] = strutil.Concat("http://", strings.TrimPrefix(cfg.NexusUrl, "http://"))
	cfgMap["NEXUS_USERNAME"] = cfg.NexusUsername
	cfgMap["NEXUS_PASSWORD"] = cfg.NexusPassword
	if err := render.RenderTemplate(gradleConfigPath, cfgMap); err != nil {
		return err
	}

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

	if _, err := os.Stat(cfg.Target); err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("target file %s not exist", cfg.Target))
		}
		return err
	}
	if err := cp(cfg.Target, cfg.WorkDir); err != nil {
		return err
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
