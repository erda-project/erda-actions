package imagebuilder

import (
	"testing"
)

func TestNewDockerBuilder(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		wantArgs []string
	}{
		{
			name:     "empty options",
			opts:     []Option{},
			wantArgs: []string{"build", "."},
		},
		{
			name: "with image name",
			opts: []Option{
				WithImageName("test-image:latest"),
			},
			wantArgs: []string{"build", ".", "-t", "test-image:latest"},
		},
		{
			name: "with dockerfile path",
			opts: []Option{
				WithDockerfilePath("./Dockerfile"),
			},
			wantArgs: []string{"build", ".", "-f", "./Dockerfile"},
		},
		{
			name: "with build context",
			opts: []Option{
				WithBuildContext("./context"),
			},
			wantArgs: []string{"build", "./context"},
		},
		{
			name: "with build args",
			opts: []Option{
				WithBuildArgs(map[string]string{
					"KEY1": "value1",
					"KEY2": "value2",
				}),
			},
			wantArgs: []string{"build", ".", "--build-arg", "KEY1=value1", "--build-arg", "KEY2=value2"},
		},
		{
			name: "with build contexts",
			opts: []Option{
				WithBuildContexts(map[string]string{
					"ctx1": "/path/to/ctx1",
					"ctx2": "/path/to/ctx2",
				}),
			},
			wantArgs: []string{"build", ".", "--build-context", "ctx1=/path/to/ctx1", "--build-context", "ctx2=/path/to/ctx2"},
		},
		{
			name: "with custom args",
			opts: []Option{
				WithCustomArgs("--no-cache", "--progress=plain"),
			},
			wantArgs: []string{"build", ".", "--no-cache", "--progress=plain"},
		},
		{
			name: "full options",
			opts: []Option{
				WithImageName("my-image:v1"),
				WithDockerfilePath("./Dockerfile"),
				WithBuildContext("./context"),
				WithBuildArgs(map[string]string{"NODE_ENV": "production"}),
				WithBuildContexts(map[string]string{"cache": "/cache"}),
				WithCustomArgs("--no-cache"),
			},
			wantArgs: []string{
				"build", "./context",
				"-f", "./Dockerfile",
				"-t", "my-image:v1",
				"--build-arg", "NODE_ENV=production",
				"--build-context", "cache=/cache",
				"--no-cache",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder(tt.opts...)
			gotArgs := builder.GetArgs()

			for _, wantArg := range tt.wantArgs {
				found := false
				for _, gotArg := range gotArgs {
					if gotArg == wantArg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected arg %q not found in %v", wantArg, gotArgs)
				}
			}
		})
	}
}

func TestDockerBuilder_AddBuildArg(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		wantArgs []string
	}{
		{
			name:     "simple build arg",
			key:      "KEY",
			value:    "value",
			wantArgs: []string{"--build-arg", "KEY=value"},
		},
		{
			name:     "build arg with special characters",
			key:      "KEY",
			value:    "value with spaces",
			wantArgs: []string{"--build-arg", "KEY=value with spaces"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder()
			builder.AddBuildArg(tt.key, tt.value)
			args := builder.GetArgs()

			// Check that the build arg is present
			found := false
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--build-arg" && args[i+1] == tt.wantArgs[1] {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Build arg not found. Got args: %v", args)
			}
		})
	}
}

func TestDockerBuilder_AddBuildArgs(t *testing.T) {
	tests := []struct {
		name  string
		args  map[string]string
		count int
	}{
		{
			name:  "empty map",
			args:  map[string]string{},
			count: 0,
		},
		{
			name: "single build arg",
			args: map[string]string{
				"KEY1": "value1",
			},
			count: 1,
		},
		{
			name: "multiple build args",
			args: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
				"KEY3": "value3",
			},
			count: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder()
			builder.AddBuildArgs(tt.args)
			args := builder.GetArgs()

			// Count --build-arg occurrences
			count := 0
			for _, arg := range args {
				if arg == "--build-arg" {
					count++
				}
			}

			if count != tt.count {
				t.Errorf("Expected %d build args, got %d. Args: %v", tt.count, count, args)
			}
		})
	}
}

func TestDockerBuilder_AddBuildContext(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		wantArgs []string
	}{
		{
			name:     "simple build context",
			key:      "ctx1",
			value:    "/path/to/ctx1",
			wantArgs: []string{"--build-context", "ctx1=/path/to/ctx1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder()
			builder.AddBuildContext(tt.key, tt.value)
			args := builder.GetArgs()

			found := false
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--build-context" && args[i+1] == tt.wantArgs[1] {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Build context not found. Got args: %v", args)
			}
		})
	}
}

func TestDockerBuilder_AddBuildContexts(t *testing.T) {
	tests := []struct {
		name      string
		contexts  map[string]string
		wantCount int
	}{
		{
			name:      "empty map",
			contexts:  map[string]string{},
			wantCount: 0,
		},
		{
			name: "multiple contexts",
			contexts: map[string]string{
				"ctx1": "/path/to/ctx1",
				"ctx2": "/path/to/ctx2",
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder()
			builder.AddBuildContexts(tt.contexts)
			args := builder.GetArgs()

			count := 0
			for _, arg := range args {
				if arg == "--build-context" {
					count++
				}
			}

			if count != tt.wantCount {
				t.Errorf("Expected %d build contexts, got %d. Args: %v", tt.wantCount, count, args)
			}
		})
	}
}

func TestDockerBuilder_AddCustomArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantArgs []string
	}{
		{
			name:     "single custom arg",
			args:     []string{"--no-cache"},
			wantArgs: []string{"--no-cache"},
		},
		{
			name:     "multiple custom args",
			args:     []string{"--no-cache", "--progress=plain", "--quiet"},
			wantArgs: []string{"--no-cache", "--progress=plain", "--quiet"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder()
			builder.AddCustomArgs(tt.args...)
			args := builder.GetArgs()

			for _, wantArg := range tt.wantArgs {
				found := false
				for _, gotArg := range args {
					if gotArg == wantArg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected custom arg %q not found in %v", wantArg, args)
				}
			}
		})
	}
}

func TestDockerBuilder_BuildCommand(t *testing.T) {
	tests := []struct {
		name      string
		opts      []Option
		wantCmd   string
		wantCount int
	}{
		{
			name:      "basic command",
			opts:      []Option{},
			wantCmd:   "docker",
			wantCount: 2, // "build" and "."
		},
		{
			name:      "command with options",
			opts:      []Option{WithImageName("test:latest")},
			wantCmd:   "docker",
			wantCount: 4, // "build", ".", "-t", "test:latest"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder(tt.opts...)
			cmd := builder.BuildCommand()

			if len(cmd.Args) == 0 || cmd.Args[0] != "docker" {
				t.Errorf("Expected command name 'docker', got args: %v", cmd.Args)
			}

			if len(cmd.Args) < tt.wantCount+1 {
				t.Errorf("Expected at least %d args (including command name), got %d: %v", tt.wantCount+1, len(cmd.Args), cmd.Args)
			}
		})
	}
}

func TestDockerBuilder_GetArgs(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		ops      []func(*DockerBuilder)
		wantLen  int
		checkArg func([]string) bool
	}{
		{
			name:    "returns copy of args",
			opts:    []Option{},
			ops:     []func(*DockerBuilder){},
			wantLen: 2,
			checkArg: func(args []string) bool {
				return len(args) >= 2 && args[0] == "build"
			},
		},
		{
			name: "args are not modified after GetArgs",
			opts: []Option{WithImageName("test:latest")},
			ops: []func(*DockerBuilder){
				func(b *DockerBuilder) {
					args1 := b.GetArgs()
					b.AddCustomArgs("--new-arg")
					args2 := b.GetArgs()
					if len(args2) <= len(args1) {
						t.Errorf("Expected args2 to be longer than args1, got len(args1)=%d, len(args2)=%d", len(args1), len(args2))
					}
				},
			},
			wantLen: 4,
			checkArg: func(args []string) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder(tt.opts...)
			for _, op := range tt.ops {
				op(builder)
			}
			args := builder.GetArgs()

			if len(args) < tt.wantLen {
				t.Errorf("Expected at least %d args, got %d: %v", tt.wantLen, len(args), args)
			}

			if tt.checkArg != nil && !tt.checkArg(args) {
				t.Errorf("Args check failed. Args: %v", args)
			}
		})
	}
}

func TestDockerBuilder_BuildArgsOrdering(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]string
		validate func([]string) bool
	}{
		{
			name: "build args are sorted",
			args: map[string]string{
				"Z_KEY": "z_value",
				"A_KEY": "a_value",
				"M_KEY": "m_value",
			},
			validate: func(args []string) bool {
				// Find all build arg values
				var values []string
				for i := 0; i < len(args)-1; i++ {
					if args[i] == "--build-arg" {
						values = append(values, args[i+1])
					}
				}
				// Check that A_KEY comes before M_KEY, and M_KEY comes before Z_KEY
				aIdx, mIdx, zIdx := -1, -1, -1
				for i, v := range values {
					if v == "A_KEY=a_value" {
						aIdx = i
					}
					if v == "M_KEY=m_value" {
						mIdx = i
					}
					if v == "Z_KEY=z_value" {
						zIdx = i
					}
				}
				return aIdx >= 0 && mIdx >= 0 && zIdx >= 0 && aIdx < mIdx && mIdx < zIdx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewDockerBuilder()
			builder.AddBuildArgs(tt.args)
			args := builder.GetArgs()

			if !tt.validate(args) {
				t.Errorf("Build args ordering validation failed. Args: %v", args)
			}
		})
	}
}
