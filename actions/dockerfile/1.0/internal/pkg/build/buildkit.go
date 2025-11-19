package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
	pkgconf "github.com/erda-project/erda-actions/pkg/envconf"
	"github.com/erda-project/erda-infra/pkg/strutil"
)

const (
	defaultCaCertPath  = "/.buildkit/ca.pem"
	defaultCertPath    = "/.buildkit/cert.pem"
	defaultCertKeyPath = "/.buildkit/key.pem"
)

type buildkit struct {
	config *conf.Conf
	args   []string
}

func NewBuildkit(c *conf.Conf) (Builder, error) {
	dfDir, dfFileName, err := resolveDockerfilePath(c)
	if err != nil {
		return nil, err
	}

	args := []string{
		"--addr", c.BuildkitdAddr,
		fmt.Sprintf("--tlscacert=%s", defaultCaCertPath),
		fmt.Sprintf("--tlscert=%s", defaultCertPath),
		fmt.Sprintf("--tlskey=%s", defaultCertKeyPath),
		"build",
		"--frontend", "dockerfile.v0",
		"--local", "context=" + c.Context,
		"--local", "dockerfile=" + dfDir,
		"--opt", "filename=" + dfFileName,
		"--opt", fmt.Sprintf("platform=%s", pkgconf.GetTargetPlatforms()),
	}

	return &buildkit{
		config: c,
		args:   args,
	}, nil
}

func (b *buildkit) Build(p *Params, o *OutPut) error {
	// generate command from action params provided
	b.appendBuildArgs(p.Args)
	b.appendBuildContexts(p.BuildContext)

	// set output
	b.args = append(b.args,
		"--output",
		fmt.Sprintf("type=image,name=%s,push=%s", o.Image, strconv.FormatBool(o.Push)),
	)

	buildCmd := exec.Command("buildctl", b.args...)
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

func resolveDockerfilePath(c *conf.Conf) (string, string, error) {
	p := c.Path
	if !filepath.IsAbs(p) {
		p = filepath.Join(c.Context, p)
	}

	fi, err := os.Stat(p)
	if err != nil {
		return "", "", err
	}

	if fi.IsDir() {
		return p, "", nil
	}

	return filepath.Dir(p), filepath.Base(p), nil
}

func (b *buildkit) appendBuildArgs(args map[string]string) {
	for k, v := range args {
		b.args = append(b.args, "--opt", fmt.Sprintf("build-arg:%s=%s", k, v))
	}
}

func (b *buildkit) appendBuildContexts(buildContexts map[string]string) {
	for k, v := range buildContexts {
		b.args = append(b.args,
			"--local", fmt.Sprintf("%s=%s", k, v),
			"--opt", fmt.Sprintf("context:%s=local:%s", k, k))
	}
}
