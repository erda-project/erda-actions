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

package repo

import (
	"path/filepath"

	"github.com/erda-project/erda-actions/actions/archive-release/1.0/internal/config"
)

type LintConfig struct {
	conf *config.Config
}

func (l *LintConfig) Bucket() string {
	return l.conf.OssBucket
}

// Remote is like /archived-versions/{git-tag:v1.0.0}/sqls/config.yml
func (l *LintConfig) Remote() string {
	return filepath.Join(l.conf.GetOssPath(), "sqls", "config.yml")
}

func (l *LintConfig) Local() string {
	if l.conf.Workdir == "" || l.conf.LintConfig == "" {
		return ""
	}
	return filepath.Join(l.conf.Workdir, l.conf.LintConfig)
}
