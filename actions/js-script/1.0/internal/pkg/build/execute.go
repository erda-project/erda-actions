package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/js-script/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

var (
	scriptPath = "/opt/action/command.sh"
)

func Execute() error {

	var cfg conf.Conf
	envconf.MustLoad(&cfg)
	if err := os.Chdir(cfg.Context); err != nil {
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

	for _, target := range cfg.Targets {
		if _, err := os.Stat(target); err != nil {
			if os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("target file %s not exist", target))
			}
			return err
		}
		if err := cp(target, cfg.WorkDir); err != nil {
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

func cp(srcDir, destDir string) error {
	fmt.Fprintf(os.Stdout, "copying %s to %s\n", srcDir, destDir)
	cpCmd := exec.Command("cp", "-r", srcDir, destDir)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	return cpCmd.Run()
}

func setupScript(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, `$`, `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		buildScript,
		buf.String(),
	)
	return script
}

// buildScript is a helper script which add a shebang
// to the generated script.
const buildScript = `#!/bin/sh
set -e
%s
`

// traceScript is a helper script which is added to the
// generated script to trace each command.
const traceScript = `
echo + %s
%s || ((echo "- FAIL! exit code: $?") && false)
echo
`
