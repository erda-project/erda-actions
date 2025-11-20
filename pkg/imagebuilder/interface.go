package imagebuilder

import "os/exec"

// Builder defines the interface for building image build parameters
type Builder interface {
	// BuildCommand builds the command with all configured parameters
	BuildCommand() *exec.Cmd
	// AddBuildArg adds a build argument
	AddBuildArg(key, value string)
	// AddBuildArgs adds multiple build arguments
	AddBuildArgs(args map[string]string)
	// AddBuildContext adds a build context mapping
	AddBuildContext(key, value string)
	// AddBuildContexts adds multiple build context mappings
	AddBuildContexts(contexts map[string]string)
	// AddCustomArgs adds custom arguments
	AddCustomArgs(args ...string)
	// GetArgs returns the built arguments
	GetArgs() []string
}

// BuildOptions contains configuration options for image builders
type BuildOptions struct {
	// ImageName is the name of the image to build
	ImageName string
	// DockerfilePath is the path to the Dockerfile
	DockerfilePath string
	// BuildContext is the build context directory
	BuildContext string
	// Push indicates whether to push the image after building
	Push bool
	// CustomArgs are additional custom arguments
	CustomArgs []string
	// BuildArgs are build arguments
	BuildArgs map[string]string
	// BuildContexts are build context mappings
	BuildContexts map[string]string
}

// Option is a function that modifies BuildOptions
type Option func(*BuildOptions)

// WithImageName sets the image name
func WithImageName(name string) Option {
	return func(o *BuildOptions) {
		o.ImageName = name
	}
}

// WithDockerfilePath sets the Dockerfile path
func WithDockerfilePath(path string) Option {
	return func(o *BuildOptions) {
		o.DockerfilePath = path
	}
}

// WithBuildContext sets the build context
func WithBuildContext(context string) Option {
	return func(o *BuildOptions) {
		o.BuildContext = context
	}
}

// WithPush sets whether to push the image
func WithPush(push bool) Option {
	return func(o *BuildOptions) {
		o.Push = push
	}
}

// WithCustomArgs sets custom arguments
func WithCustomArgs(args ...string) Option {
	return func(o *BuildOptions) {
		o.CustomArgs = append(o.CustomArgs, args...)
	}
}

// WithBuildArgs sets build arguments
func WithBuildArgs(args map[string]string) Option {
	return func(o *BuildOptions) {
		if o.BuildArgs == nil {
			o.BuildArgs = make(map[string]string)
		}
		for k, v := range args {
			o.BuildArgs[k] = v
		}
	}
}

// WithBuildContexts sets build context mappings
func WithBuildContexts(contexts map[string]string) Option {
	return func(o *BuildOptions) {
		if o.BuildContexts == nil {
			o.BuildContexts = make(map[string]string)
		}
		for k, v := range contexts {
			o.BuildContexts[k] = v
		}
	}
}
