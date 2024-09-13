package build

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda/apistructs"
)

type dockerBuilder struct {
	config *conf.Conf
	args   []string
}

func NewDocker(c *conf.Conf) Builder {
	args := []string{
		"build",
		".", // set current dir is build context
		"-f", c.Path,
	}

	return &dockerBuilder{
		config: c,
		args:   args,
	}
}

func (d *dockerBuilder) Build(p *Params, o *OutPut) error {
	// generate command from action params provided
	d.appendBuildArgs(p.Args)
	d.appendBuildContexts(p.BuildContext)

	// set resource limit
	d.args = append(d.args,
		"--cpu-quota", strconv.FormatFloat(d.config.CPU*100000, 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(d.config.Memory*apistructs.MB), 10),
	)

	// generate build command
	buildCmd := exec.Command("docker", d.args...)
	fmt.Fprintf(os.Stdout, "Docker build command: %v\n", buildCmd.Args)

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Dir = d.config.WorkDir

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build docker image: %v", err)
	}

	if o.Push {
		if err := docker.PushByCmd(o.Image, ""); err != nil {
			return fmt.Errorf("failed to push docker image: %v", err)
		}
	}

	fmt.Fprintf(os.Stdout, "Successfully built and pushed image: %s\n", o.Image)
	return nil
}

// appendBuildArgs
func (d *dockerBuilder) appendBuildArgs(args map[string]string) {
	for k, v := range args {
		d.args = append(d.args, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}
}

// appendBuildContexts
func (d *dockerBuilder) appendBuildContexts(contexts map[string]string) {
	for k, v := range contexts {
		d.args = append(d.args, "--build-context", fmt.Sprintf("%s=%s", k, v))
	}
}
