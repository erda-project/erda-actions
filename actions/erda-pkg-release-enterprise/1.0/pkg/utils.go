package pkg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			if _, err = w.Write(d); err != nil {
				return nil, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// ExecCmd execute shell commands
func ExecCmd(logFile, errFile *os.File, dir string, name string, args ...string) (*bytes.Buffer, error) {

	var outPut, outErr []byte
	var stdErr error

	command := fmt.Sprintf("%s %s", name, strings.Join(args, " "))
	WriteLog(logFile, fmt.Sprintf("exec command: %s\n", command))

	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		WriteLog(errFile, fmt.Sprintf("start command: %s failed. err: %v", command, err))
		return nil, err
	}

	go func() {
		outPut, stdErr = copyAndCapture(logFile, stdoutIn)
	}()
	go func() {
		outErr, stdErr = copyAndCapture(errFile, stderrIn)
	}()

	if err := cmd.Wait(); err != nil {
		WriteLog(errFile, fmt.Sprintf("wait command: %s failed. err: %v", command, err))
		return nil, err
	}

	if stdErr != nil && strings.Index(stdErr.Error(), "closed") < 0 {
		return nil, errors.New(bytes.NewBuffer(outErr).String())
	}

	return bytes.NewBuffer(outPut), nil
}

// WriteLog write log
func WriteLog(logFile *os.File, info string) {
	_, _ = logFile.WriteString(info)
}

// IsDirExists judge dir exists or not
func IsDirExists(path string) (bool, error) {
	si, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	if !si.IsDir() {
		return false, fmt.Errorf("source: %s is not a directory", path)
	}

	return true, nil
}

// FileExist judge file exists or not
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GenErdaPublicName
func GenErdaPublicName(version, osArch string) string {
	return fmt.Sprintf("erda-%s.%s.tar.gz", version, osArch)
}
