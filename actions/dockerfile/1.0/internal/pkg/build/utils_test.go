package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
)

func Test_resolveDockerfilePath(t *testing.T) {
	tmp := t.TempDir()

	existingDir := filepath.Join(tmp, "project")
	existingFile := filepath.Join(tmp, "project", "Dockerfile")
	existingDevFile := filepath.Join(tmp, "project", "Dockerfile.dev")

	os.MkdirAll(existingDir, 0o755)
	os.WriteFile(existingFile, []byte("FROM busybox"), 0o644)
	os.WriteFile(existingDevFile, []byte("FROM busybox"), 0o644)

	tests := []struct {
		name     string
		conf     *conf.Conf
		wantDir  string
		wantFile string
		wantErr  bool
	}{
		{
			name: "relative path + context",
			conf: &conf.Conf{
				Context: tmp,
				Path:    "project/Dockerfile",
			},
			wantDir:  filepath.Join(tmp, "project"),
			wantFile: "Dockerfile",
			wantErr:  false,
		},
		{
			name: "absolute file path",
			conf: &conf.Conf{
				Path: existingFile,
			},
			wantDir:  filepath.Dir(existingFile),
			wantFile: filepath.Base(existingFile),
			wantErr:  false,
		},
		{
			name: "directory path",
			conf: &conf.Conf{
				Context: tmp,
				Path:    "project",
			},
			wantDir:  filepath.Join(tmp, "project"),
			wantFile: "",
			wantErr:  false,
		},
		{
			name: "non-existing path",
			conf: &conf.Conf{
				Context: tmp,
				Path:    "not-exist",
			},
			wantErr: true,
		},
		{
			name: "empty context but absolute file path",
			conf: &conf.Conf{
				Path: existingFile,
			},
			wantDir:  filepath.Dir(existingFile),
			wantFile: filepath.Base(existingFile),
			wantErr:  false,
		},
		{
			name: "empty context + relative path should error",
			conf: &conf.Conf{
				Context: "",
				Path:    "Dockerfile",
			},
			wantErr: true,
		},
		{
			name: "custom dockerfile",
			conf: &conf.Conf{
				Context: tmp,
				Path:    "project/Dockerfile.dev",
			},
			wantDir:  filepath.Join(tmp, "project"),
			wantFile: "Dockerfile.dev",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, file, err := resolveDockerfilePath(tt.conf)

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v, wantErr=%v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if dir != tt.wantDir {
				t.Fatalf("dir mismatch: got %s, want %s", dir, tt.wantDir)
			}

			if file != tt.wantFile {
				t.Fatalf("file mismatch: got %s, want %s", file, tt.wantFile)
			}
		})
	}
}
