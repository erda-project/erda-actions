package pack

import (
	"os"
	"os/exec"
)

func Tar(tarAbsPath string, sourcePath string) error {
	args := []string{"-cf", tarAbsPath, sourcePath}
	cmd := exec.Command("tar", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func UnTar(tarAbsPath string, destDir string) error {
	os.Mkdir(destDir, os.ModePerm)
	cmd := exec.Command("tar", "-xf", tarAbsPath, "-C", destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
