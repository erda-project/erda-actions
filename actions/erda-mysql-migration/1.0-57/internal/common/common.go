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

package common

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	mysqld = "/usr/bin/run-mysqld"
)

// StartSandbox start a MySQL server in the container
func StartSandbox() error {
	logrus.Infoln("Create sandbox")
	envs := os.Environ()
	sandbox := exec.Command(mysqld)
	for _, env := range envs {
		for _, key := range []string{"MYSQL_USER", "MYSQL_USERNAME", "MYSQL_PASSWORD", "MYSQL_DATABASE", "MYSQL_HOST", "MYSQL_PORT"} {
			if !strings.HasPrefix(env, key+"=") {
				sandbox.Env = append(sandbox.Env, env)
			}
		}
	}
	if err := sandbox.Start(); err != nil {
		return errors.Wrap(err, "failed to Start sandbox")
	}
	if err := sandbox.Wait(); err != nil {
		return errors.Wrapf(err, "failed to exec %s", mysqld)
	}
	return nil
}

func FatalError(f func() error) {
	if err := f(); err != nil {
		logrus.Fatalln(err)
	}
}
