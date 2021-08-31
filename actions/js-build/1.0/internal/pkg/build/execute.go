package build

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/erda-project/erda-actions/actions/js-build/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

var SupportNodeVersionFormatMap = map[string]string{
	"8":  "v8.17.0",
	"10": "v10.24.1",
	"12": "v12.22.5",
	"14": "v14.17.5",
}

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

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

	userVersion := SupportNodeVersionFormatMap[cfg.NodeVersion]

	if userVersion == "" {
		return errors.New(fmt.Sprintf(" not support this node version %v ", cfg.NodeVersion))
	}

	var cmdStr string
	cmdStr += "source ~/.bashrc"
	cmdStr += " && nvm use " + userVersion

	cmd := cfg.BuildCmd
	if cmd == nil || len(cmd) == 0 {
		return errors.New(" error BuildCmd is empty ")
	}

	for _, cmd := range cfg.BuildCmd {
		cmdStr += " && " + cmd
	}

	if err := runCommand(cmdStr); err != nil {
		return err
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
	runCommand(fmt.Sprintf(" cp -r %s/. ./", cfg.Context))

	//创建输出OUTPUT文件
	if !filepath.IsAbs(cfg.MetaFile) {
		return errors.New(fmt.Sprintf("not an absolute path: %s", cfg.MetaFile))
	}
	err := os.MkdirAll(filepath.Dir(cfg.MetaFile), 0755)
	if err != nil {
		return err
	}

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
