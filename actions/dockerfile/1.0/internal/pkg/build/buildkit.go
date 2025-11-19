package build

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
	pkgconf "github.com/erda-project/erda-actions/pkg/envconf"
	"github.com/erda-project/erda-actions/pkg/imagebuilder"
	"github.com/erda-project/erda-infra/pkg/strutil"
)

type buildkit struct {
	config         *conf.Conf
	dockerfileDir  string
	dockerfileName string
}

func NewBuildkit(c *conf.Conf) (Builder, error) {
	dfDir, dfFileName, err := resolveDockerfilePath(c)
	if err != nil {
		return nil, err
	}

	return &buildkit{
		config:         c,
		dockerfileDir:  dfDir,
		dockerfileName: dfFileName,
	}, nil
}

func (b *buildkit) Build(p *Params, o *OutPut) error {
	buildkitOpts := &imagebuilder.BuildKitOptions{
		Addr:           b.config.BuildkitdAddr,
		Platform:       pkgconf.GetTargetPlatforms(),
		DockerfileDir:  b.dockerfileDir,
		DockerfileName: b.dockerfileName,
	}

	builder := imagebuilder.NewBuildKitBuilder(
		buildkitOpts,
		imagebuilder.WithImageName(o.Image),
		imagebuilder.WithBuildContext(b.config.Context),
		imagebuilder.WithPush(o.Push),
		imagebuilder.WithBuildArgs(p.Args),
		imagebuilder.WithBuildContexts(p.BuildContext),
	)

	buildCmd := builder.BuildCommand()
	fmt.Println(strutil.Join(buildCmd.Args, " ", false))
	buildCmd.Dir = b.config.WorkDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build image: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully built and pushed image: %s\n", o.Image)
	return nil
}
