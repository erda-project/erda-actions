package imagebuilder

import (
	"fmt"
	"os/exec"
	"sort"
)

const (
	// DefaultTLSCACertPath is the default path to TLS CA certificate
	DefaultTLSCACertPath = "/.buildkit/ca.pem"
	// DefaultTLSCertPath is the default path to TLS certificate
	DefaultTLSCertPath = "/.buildkit/cert.pem"
	// DefaultTLSKeyPath is the default path to TLS key
	DefaultTLSKeyPath = "/.buildkit/key.pem"
	// DefaultFrontend is the default BuildKit frontend
	DefaultFrontend = "dockerfile.v0"
)

// BuildKitOptions contains BuildKit-specific configuration
type BuildKitOptions struct {
	// Addr is the BuildKit daemon address
	Addr string
	// TLSCACertPath is the path to TLS CA certificate, defaults to DefaultTLSCACertPath
	TLSCACertPath string
	// TLSCertPath is the path to TLS certificate, defaults to DefaultTLSCertPath
	TLSCertPath string
	// TLSKeyPath is the path to TLS key, defaults to DefaultTLSKeyPath
	TLSKeyPath string
	// Frontend is the BuildKit frontend to use, defaults to DefaultFrontend
	Frontend string
	// Platform is the target platform
	Platform string
	// DockerfileDir is the directory containing the Dockerfile
	DockerfileDir string
	// DockerfileName is the name of the Dockerfile
	DockerfileName string
}

// BuildKitBuilder implements Builder interface for BuildKit CLI
type BuildKitBuilder struct {
	options      *BuildOptions
	buildkitOpts *BuildKitOptions
	args         []string
}

// NewBuildKitBuilder creates a new BuildKit builder
func NewBuildKitBuilder(buildkitOpts *BuildKitOptions, opts ...Option) *BuildKitBuilder {
	options := &BuildOptions{
		BuildArgs:     make(map[string]string),
		BuildContexts: make(map[string]string),
		CustomArgs:    make([]string, 0),
	}
	for _, opt := range opts {
		opt(options)
	}

	if buildkitOpts == nil {
		buildkitOpts = &BuildKitOptions{}
	}

	// Apply default values for TLS certificates if not provided
	if buildkitOpts.TLSCACertPath == "" {
		buildkitOpts.TLSCACertPath = DefaultTLSCACertPath
	}
	if buildkitOpts.TLSCertPath == "" {
		buildkitOpts.TLSCertPath = DefaultTLSCertPath
	}
	if buildkitOpts.TLSKeyPath == "" {
		buildkitOpts.TLSKeyPath = DefaultTLSKeyPath
	}
	if buildkitOpts.Frontend == "" {
		buildkitOpts.Frontend = DefaultFrontend
	}

	builder := &BuildKitBuilder{
		options:      options,
		buildkitOpts: buildkitOpts,
		args:         make([]string, 0),
	}

	builder.buildBaseArgs()
	builder.AddBuildArgs(options.BuildArgs)
	builder.AddBuildContexts(options.BuildContexts)
	if len(options.CustomArgs) > 0 {
		builder.AddCustomArgs(options.CustomArgs...)
	}
	return builder
}

func (b *BuildKitBuilder) buildBaseArgs() {
	if b.buildkitOpts.Addr != "" {
		b.args = append(b.args, "--addr", b.buildkitOpts.Addr)
	}
	// Always add TLS certificates with default values
	b.args = append(b.args, fmt.Sprintf("--tlscacert=%s", b.buildkitOpts.TLSCACertPath))
	b.args = append(b.args, fmt.Sprintf("--tlscert=%s", b.buildkitOpts.TLSCertPath))
	b.args = append(b.args, fmt.Sprintf("--tlskey=%s", b.buildkitOpts.TLSKeyPath))

	b.args = append(b.args, "build")
	b.args = append(b.args, "--frontend", b.buildkitOpts.Frontend)

	if b.options.BuildContext != "" {
		b.args = append(b.args, "--local", fmt.Sprintf("context=%s", b.options.BuildContext))
	}

	if b.buildkitOpts.DockerfileDir != "" {
		b.args = append(b.args, "--local", fmt.Sprintf("dockerfile=%s", b.buildkitOpts.DockerfileDir))
	}

	if b.buildkitOpts.DockerfileName != "" {
		b.args = append(b.args, "--opt", fmt.Sprintf("filename=%s", b.buildkitOpts.DockerfileName))
	}

	if b.buildkitOpts.Platform != "" {
		b.args = append(b.args, "--opt", fmt.Sprintf("platform=%s", b.buildkitOpts.Platform))
	}

	if b.options.ImageName != "" {
		output := fmt.Sprintf("type=image,name=%s", b.options.ImageName)
		if b.options.Push {
			output += ",push=true"
		} else {
			output += ",push=false"
		}
		b.args = append(b.args, "--output", output)
	}
}

// BuildCommand builds the buildctl build command
func (b *BuildKitBuilder) BuildCommand() *exec.Cmd {
	return exec.Command("buildctl", b.args...)
}

// AddBuildArg adds a build argument
func (b *BuildKitBuilder) AddBuildArg(key, value string) {
	b.args = append(b.args, "--opt", fmt.Sprintf("build-arg:%s=%s", key, value))
}

// AddBuildArgs adds multiple build arguments
func (b *BuildKitBuilder) AddBuildArgs(args map[string]string) {
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
func (b *BuildKitBuilder) AddBuildContext(key, value string) {
	b.args = append(b.args, "--local", fmt.Sprintf("%s=%s", key, value))
	b.args = append(b.args, "--opt", fmt.Sprintf("context:%s=local:%s", key, key))
}

// AddBuildContexts adds multiple build context mappings
func (b *BuildKitBuilder) AddBuildContexts(contexts map[string]string) {
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
func (b *BuildKitBuilder) AddCustomArgs(args ...string) {
	b.args = append(b.args, args...)
}

// GetArgs returns the built arguments
func (b *BuildKitBuilder) GetArgs() []string {
	return append([]string(nil), b.args...)
}
