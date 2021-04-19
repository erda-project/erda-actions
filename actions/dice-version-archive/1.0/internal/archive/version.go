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

// read VERSIION file from dice repo
package archive

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type Version struct {
	version string
	major   uint64
	minor   uint64
}

func (v *Version) Read(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data = bytes.TrimFunc(data, unicode.IsSpace)
	v.version = string(data)
	split := strings.Split(v.version, ".")
	if len(split) != 2 {
		return errors.New("can not parse version from VERSION file")
	}
	major, err := strconv.ParseUint(split[0], 10, 32)
	if err != nil {
		return errors.New("can not parse major version from VERSION file")
	}
	minor, err := strconv.ParseUint(split[1], 10, 32)
	if err != nil {
		return errors.New("can not parse minor version from VERSION file")
	}
	v.major = major
	v.minor = minor

	return nil
}

func (v *Version) String() string {
	return v.version
}

func (v *Version) Major() uint64 {
	return v.major
}

func (v *Version) Minor() uint64 {
	return v.minor
}
