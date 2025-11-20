package build

import (
	"os"
	"path/filepath"

	"github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/pkg/conf"
)

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
