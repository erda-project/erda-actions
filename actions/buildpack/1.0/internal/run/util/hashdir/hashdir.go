package hashdir

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// SHA1, SHA256, MD5
const (
	SHA1   = "sha1"
	SHA256 = "sha256"
	MD5    = "md5"
)

// GetHash create an instance of specific hash algorithm
func GetHash(name *string) (hash.Hash, error) {
	if *name == SHA1 {
		return sha1.New(), nil
	} else if *name == SHA256 {
		return sha256.New(), nil
	} else if *name == MD5 {
		return md5.New(), nil
	}
	message := "Hash Algorithm is not supported"
	err := errors.New(message)
	return nil, err

}

// Create hash value with local path and a hash algorithm
func Create(dir string, hashAlgorithm string, ignorePatterns []string) (string, error) {

	h, err := GetHash(&hashAlgorithm)

	if err != nil {
		return "", nil
	}

	gitignore := ignore.CompileIgnoreLines(ignorePatterns...)

	err = filepath.Walk(dir, func(absPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		relativePath := filepath.Clean(strings.TrimPrefix(absPath, filepath.Clean(dir)+"/"))
		if gitignore.MatchesPath(relativePath) {
			return nil
		}
		b, err := ioutil.ReadFile(absPath)
		if err != nil {
			return nil
		}
		io.WriteString(h, absPath+string(b))
		return nil
	})

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
