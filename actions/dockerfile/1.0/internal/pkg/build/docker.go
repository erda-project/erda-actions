package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/erda-project/erda/apistructs"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda-actions/pkg/imagebuilder"
)

type dockerBuilder struct {
	config         *conf.Conf
	dockerfilePath string
}

func NewDocker(c *conf.Conf) (Builder, error) {
	dfDir, dirFileName, err := resolveDockerfilePath(c)
	if err != nil {
		return nil, err
	}

	return &dockerBuilder{
		config:         c,
		dockerfilePath: filepath.Join(dfDir, dirFileName),
	}, nil
}

func (d *dockerBuilder) Build(p *Params, o *OutPut) error {
	buildContext := d.config.Context
	if buildContext == "" {
		buildContext = "."
	}

	builder := imagebuilder.NewDockerBuilder(
		imagebuilder.WithImageName(o.Image),
		imagebuilder.WithDockerfilePath(d.dockerfilePath),
		imagebuilder.WithBuildContext(buildContext),
		imagebuilder.WithBuildArgs(p.Args),
		imagebuilder.WithBuildContexts(p.BuildContext),
	)

	// set resource limit
	builder.AddCustomArgs(
		"--cpu-quota", strconv.FormatFloat(d.config.CPU*100000, 'f', 0, 64),
		"--memory", strconv.FormatInt(int64(d.config.Memory*apistructs.MB), 10),
	)

	buildCmd := builder.BuildCommand()
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
