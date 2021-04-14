// 读取 dice 仓库下的 VERSIION 文件

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
