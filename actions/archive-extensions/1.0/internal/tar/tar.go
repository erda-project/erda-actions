// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package tar

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func Tar(src, dst string) (string, error) {
	src, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}
	if dst == "" {
		dst = src + ".tar.gz"
	}

	tar := exec.Command("tar", "-zcf", dst, "-C", filepath.Dir(src), filepath.Base(src))
	tar.Stdout = os.Stdout
	tar.Stderr = os.Stderr
	if err := tar.Run(); err != nil {
		return "", errors.Wrapf(err, "failed to tar.Run, command: %s", tar.String())
	}

	return dst, nil
}
