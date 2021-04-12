package utils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func Parse(names ...string) ([]byte, error) {
	for _, name := range names {
		if fileExists(name) {
			bytes, err := ioutil.ReadFile(name)
			if err != nil {
				return nil, err
			}
			return bytes, nil
		}
	}
	return nil, errors.Errorf("can not find file with filename = %v", names)
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFile(absPath, content string, perm os.FileMode) error {
	if !filepath.IsAbs(absPath) {
		return errors.Errorf("not an absolute path: %s", absPath)
	}
	err := os.MkdirAll(filepath.Dir(absPath), 0755)
	if err != nil {
		return errors.Wrap(err, "make parent dir error")
	}
	f, err := os.OpenFile(filepath.Clean(absPath), os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return errors.Wrap(err, "write content to file error")
	}
	return nil
}

func MakeSrcDir(name string) error {
	if err := os.MkdirAll(name, 0755); err != nil {
		return errors.Wrapf(err, "create directory error: %s", name)
	}
	return nil
}

func ExecScript(scriptPath string) error {
	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
