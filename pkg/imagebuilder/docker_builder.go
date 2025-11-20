package imagebuilder

import (
	"fmt"
	"os/exec"
	"sort"
)

// DockerBuilder implements Builder interface for Docker CLI
type DockerBuilder struct {
	options *BuildOptions
	args    []string
}

// NewDockerBuilder creates a new Docker builder
func NewDockerBuilder(opts ...Option) *DockerBuilder {
	options := &BuildOptions{
		BuildArgs:     make(map[string]string),
		BuildContexts: make(map[string]string),
		CustomArgs:    make([]string, 0),
	}
	for _, opt := range opts {
		opt(options)
	}

	builder := &DockerBuilder{
		options: options,
		args:    make([]string, 0),
	}

	builder.buildBaseArgs()
	builder.AddBuildArgs(options.BuildArgs)
	builder.AddBuildContexts(options.BuildContexts)
	if len(options.CustomArgs) > 0 {
		builder.AddCustomArgs(options.CustomArgs...)
	}
	return builder
}

func (b *DockerBuilder) buildBaseArgs() {
	b.args = append(b.args, "build")
	if b.options.BuildContext != "" {
		b.args = append(b.args, b.options.BuildContext)
	} else {
		b.args = append(b.args, ".")
	}
	if b.options.DockerfilePath != "" {
		b.args = append(b.args, "-f", b.options.DockerfilePath)
	}
	if b.options.ImageName != "" {
		b.args = append(b.args, "-t", b.options.ImageName)
	}
}

// BuildCommand builds the docker build command
func (b *DockerBuilder) BuildCommand() *exec.Cmd {
	return exec.Command("docker", b.args...)
}

// AddBuildArg adds a build argument
func (b *DockerBuilder) AddBuildArg(key, value string) {
	b.args = append(b.args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
}

// AddBuildArgs adds multiple build arguments
func (b *DockerBuilder) AddBuildArgs(args map[string]string) {
	if len(args) == 0 {
		return
	}
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.AddBuildArg(key, args[key])
	}
}

// AddBuildContext adds a build context mapping
func (b *DockerBuilder) AddBuildContext(key, value string) {
	b.args = append(b.args, "--build-context", fmt.Sprintf("%s=%s", key, value))
}

// AddBuildContexts adds multiple build context mappings
func (b *DockerBuilder) AddBuildContexts(contexts map[string]string) {
	if len(contexts) == 0 {
		return
	}
	keys := make([]string, 0, len(contexts))
	for k := range contexts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.AddBuildContext(key, contexts[key])
	}
}

// AddCustomArgs adds custom arguments
func (b *DockerBuilder) AddCustomArgs(args ...string) {
	b.args = append(b.args, args...)
}

// GetArgs returns the built arguments
func (b *DockerBuilder) GetArgs() []string {
	return append([]string(nil), b.args...)
}
