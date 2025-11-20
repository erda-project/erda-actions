package imagebuilder

import (
	"strings"
	"testing"
)

func TestNewBuildKitBuilder(t *testing.T) {
	tests := []struct {
		name         string
		buildkitOpts *BuildKitOptions
		opts         []Option
		wantTLS      bool
		wantFrontend string
	}{
		{
			name:         "nil buildkit options with defaults",
			buildkitOpts: nil,
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name:         "empty buildkit options with defaults",
			buildkitOpts: &BuildKitOptions{},
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name: "with custom TLS paths",
			buildkitOpts: &BuildKitOptions{
				TLSCACertPath: "/custom/ca.pem",
				TLSCertPath:   "/custom/cert.pem",
				TLSKeyPath:    "/custom/key.pem",
			},
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name: "with custom frontend",
			buildkitOpts: &BuildKitOptions{
				Frontend: "dockerfile.v1",
			},
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: "dockerfile.v1",
		},
		{
			name: "with addr",
			buildkitOpts: &BuildKitOptions{
				Addr: "tcp://localhost:1234",
			},
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name: "with platform",
			buildkitOpts: &BuildKitOptions{
				Platform: "linux/amd64",
			},
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name: "with dockerfile dir and name",
			buildkitOpts: &BuildKitOptions{
				DockerfileDir:  "/path/to/dockerfile",
				DockerfileName: "Dockerfile.custom",
			},
			opts:         []Option{},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name:         "with image name and push",
			buildkitOpts: &BuildKitOptions{},
			opts: []Option{
				WithImageName("test-image:latest"),
				WithPush(true),
			},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
		{
			name: "full options",
			buildkitOpts: &BuildKitOptions{
				Addr:           "tcp://localhost:1234",
				Platform:       "linux/amd64",
				DockerfileDir:  "/path/to/dockerfile",
				DockerfileName: "Dockerfile",
			},
			opts: []Option{
				WithImageName("my-image:v1"),
				WithBuildContext("./context"),
				WithBuildArgs(map[string]string{"KEY": "value"}),
				WithBuildContexts(map[string]string{"ctx": "/path"}),
				WithPush(true),
			},
			wantTLS:      true,
			wantFrontend: DefaultFrontend,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(tt.buildkitOpts, tt.opts...)
			args := builder.GetArgs()

			// Check TLS certificates are present
			if tt.wantTLS {
				hasTLS := false
				for _, arg := range args {
					if arg == "--tlscacert="+DefaultTLSCACertPath || arg == "--tlscert="+DefaultTLSCertPath || arg == "--tlskey="+DefaultTLSKeyPath {
						hasTLS = true
						break
					}
				}
				if !hasTLS {
					// Check if custom TLS paths are used
					if tt.buildkitOpts != nil && tt.buildkitOpts.TLSCACertPath != "" {
						hasTLS = true
					}
				}
				if !hasTLS {
					t.Errorf("Expected TLS certificates in args, got: %v", args)
				}
			}

			// Check frontend
			foundFrontend := false
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--frontend" && args[i+1] == tt.wantFrontend {
					foundFrontend = true
					break
				}
			}
			if !foundFrontend {
				t.Errorf("Expected frontend %q not found in args: %v", tt.wantFrontend, args)
			}
		})
	}
}

func TestBuildKitBuilder_AddBuildArg(t *testing.T) {
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
			wantArgs: []string{"--opt", "build-arg:KEY=value"},
		},
		{
			name:     "build arg with special characters",
			key:      "KEY",
			value:    "value with spaces",
			wantArgs: []string{"--opt", "build-arg:KEY=value with spaces"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(nil)
			builder.AddBuildArg(tt.key, tt.value)
			args := builder.GetArgs()

			found := false
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--opt" && args[i+1] == tt.wantArgs[1] {
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

func TestBuildKitBuilder_AddBuildArgs(t *testing.T) {
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
			builder := NewBuildKitBuilder(nil)
			builder.AddBuildArgs(tt.args)
			args := builder.GetArgs()

			count := 0
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--opt" && len(args[i+1]) > 10 && args[i+1][:10] == "build-arg:" {
					count++
				}
			}

			if count != tt.count {
				t.Errorf("Expected %d build args, got %d. Args: %v", tt.count, count, args)
			}
		})
	}
}

func TestBuildKitBuilder_AddBuildContext(t *testing.T) {
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
			wantArgs: []string{"--local", "ctx1=/path/to/ctx1", "--opt", "context:ctx1=local:ctx1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(nil)
			builder.AddBuildContext(tt.key, tt.value)
			args := builder.GetArgs()

			hasLocal := false
			hasOpt := false
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--local" && args[i+1] == tt.wantArgs[1] {
					hasLocal = true
				}
				if args[i] == "--opt" && args[i+1] == tt.wantArgs[3] {
					hasOpt = true
				}
			}
			if !hasLocal || !hasOpt {
				t.Errorf("Build context not found correctly. Got args: %v", args)
			}
		})
	}
}

func TestBuildKitBuilder_AddBuildContexts(t *testing.T) {
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
			builder := NewBuildKitBuilder(nil)
			builder.AddBuildContexts(tt.contexts)
			args := builder.GetArgs()

			count := 0
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "--local" && i+1 < len(args) {
					// Check if next arg matches pattern
					if len(args[i+1]) > 0 {
						count++
					}
				}
			}

			if count != tt.wantCount {
				t.Errorf("Expected %d build contexts, got %d. Args: %v", tt.wantCount, count, args)
			}
		})
	}
}

func TestBuildKitBuilder_AddCustomArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantArgs []string
	}{
		{
			name:     "single custom arg",
			args:     []string{"--opt", "key=value"},
			wantArgs: []string{"--opt", "key=value"},
		},
		{
			name:     "multiple custom args",
			args:     []string{"--opt", "key1=value1", "--opt", "key2=value2"},
			wantArgs: []string{"--opt", "key1=value1", "--opt", "key2=value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(nil)
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

func TestBuildKitBuilder_BuildCommand(t *testing.T) {
	tests := []struct {
		name         string
		buildkitOpts *BuildKitOptions
		opts         []Option
		wantCmd      string
		wantMinLen   int
	}{
		{
			name:         "basic command",
			buildkitOpts: nil,
			opts:         []Option{},
			wantCmd:      "buildctl",
			wantMinLen:   5, // TLS certs + build + frontend
		},
		{
			name: "command with options",
			buildkitOpts: &BuildKitOptions{
				Addr: "tcp://localhost:1234",
			},
			opts:       []Option{WithImageName("test:latest")},
			wantCmd:    "buildctl",
			wantMinLen: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(tt.buildkitOpts, tt.opts...)
			cmd := builder.BuildCommand()

			if cmd.Path != tt.wantCmd {
				t.Errorf("Expected command path %q, got %q", tt.wantCmd, cmd.Path)
			}

			if len(cmd.Args) < tt.wantMinLen {
				t.Errorf("Expected at least %d args, got %d: %v", tt.wantMinLen, len(cmd.Args), cmd.Args)
			}
		})
	}
}

func TestBuildKitBuilder_GetArgs(t *testing.T) {
	tests := []struct {
		name         string
		buildkitOpts *BuildKitOptions
		opts         []Option
		wantMinLen   int
		checkArg     func([]string) bool
	}{
		{
			name:         "returns copy of args",
			buildkitOpts: nil,
			opts:         []Option{},
			wantMinLen:   5,
			checkArg: func(args []string) bool {
				return len(args) >= 5
			},
		},
		{
			name:         "args are not modified after GetArgs",
			buildkitOpts: nil,
			opts:         []Option{WithImageName("test:latest")},
			wantMinLen:   6,
			checkArg: func(args []string) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(tt.buildkitOpts, tt.opts...)
			args := builder.GetArgs()

			if len(args) < tt.wantMinLen {
				t.Errorf("Expected at least %d args, got %d: %v", tt.wantMinLen, len(args), args)
			}

			if tt.checkArg != nil && !tt.checkArg(args) {
				t.Errorf("Args check failed. Args: %v", args)
			}
		})
	}
}

func TestBuildKitBuilder_DefaultTLSValues(t *testing.T) {
	tests := []struct {
		name         string
		buildkitOpts *BuildKitOptions
		wantCACert   string
		wantCert     string
		wantKey      string
	}{
		{
			name:         "nil options uses defaults",
			buildkitOpts: nil,
			wantCACert:   DefaultTLSCACertPath,
			wantCert:     DefaultTLSCertPath,
			wantKey:      DefaultTLSKeyPath,
		},
		{
			name:         "empty options uses defaults",
			buildkitOpts: &BuildKitOptions{},
			wantCACert:   DefaultTLSCACertPath,
			wantCert:     DefaultTLSCertPath,
			wantKey:      DefaultTLSKeyPath,
		},
		{
			name: "custom TLS paths",
			buildkitOpts: &BuildKitOptions{
				TLSCACertPath: "/custom/ca.pem",
				TLSCertPath:   "/custom/cert.pem",
				TLSKeyPath:    "/custom/key.pem",
			},
			wantCACert: "/custom/ca.pem",
			wantCert:   "/custom/cert.pem",
			wantKey:    "/custom/key.pem",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(tt.buildkitOpts)
			args := builder.GetArgs()

			hasCACert := false
			hasCert := false
			hasKey := false

			for _, arg := range args {
				if arg == "--tlscacert="+tt.wantCACert {
					hasCACert = true
				}
				if arg == "--tlscert="+tt.wantCert {
					hasCert = true
				}
				if arg == "--tlskey="+tt.wantKey {
					hasKey = true
				}
			}

			if !hasCACert {
				t.Errorf("Expected TLS CA cert path %q not found in args: %v", tt.wantCACert, args)
			}
			if !hasCert {
				t.Errorf("Expected TLS cert path %q not found in args: %v", tt.wantCert, args)
			}
			if !hasKey {
				t.Errorf("Expected TLS key path %q not found in args: %v", tt.wantKey, args)
			}
		})
	}
}

func TestBuildKitBuilder_OutputWithPush(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		wantPush bool
	}{
		{
			name:     "push true",
			opts:     []Option{WithImageName("test:latest"), WithPush(true)},
			wantPush: true,
		},
		{
			name:     "push false",
			opts:     []Option{WithImageName("test:latest"), WithPush(false)},
			wantPush: false,
		},
		{
			name:     "push default (false)",
			opts:     []Option{WithImageName("test:latest")},
			wantPush: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(nil, tt.opts...)
			args := builder.GetArgs()

			foundOutput := false
			hasPushTrue := false
			hasPushFalse := false

			for i, arg := range args {
				if arg == "--output" && i+1 < len(args) {
					foundOutput = true
					outputValue := args[i+1]
					if strings.Contains(outputValue, ",push=true") {
						hasPushTrue = true
					}
					if strings.Contains(outputValue, ",push=false") {
						hasPushFalse = true
					}
				}
			}

			if !foundOutput {
				t.Errorf("Expected --output argument not found in args: %v", args)
			}

			if tt.wantPush && !hasPushTrue {
				t.Errorf("Expected push=true in output, got args: %v", args)
			}
			if !tt.wantPush && !hasPushFalse {
				t.Errorf("Expected push=false in output, got args: %v", args)
			}
		})
	}
}

func TestBuildKitBuilder_BuildArgsOrdering(t *testing.T) {
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
				var values []string
				for i := 0; i < len(args)-1; i++ {
					if args[i] == "--opt" && len(args[i+1]) > 10 && args[i+1][:10] == "build-arg:" {
						values = append(values, args[i+1])
					}
				}
				aIdx, mIdx, zIdx := -1, -1, -1
				for i, v := range values {
					if v == "build-arg:A_KEY=a_value" {
						aIdx = i
					}
					if v == "build-arg:M_KEY=m_value" {
						mIdx = i
					}
					if v == "build-arg:Z_KEY=z_value" {
						zIdx = i
					}
				}
				return aIdx >= 0 && mIdx >= 0 && zIdx >= 0 && aIdx < mIdx && mIdx < zIdx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuildKitBuilder(nil)
			builder.AddBuildArgs(tt.args)
			args := builder.GetArgs()

			if !tt.validate(args) {
				t.Errorf("Build args ordering validation failed. Args: %v", args)
			}
		})
	}
}
