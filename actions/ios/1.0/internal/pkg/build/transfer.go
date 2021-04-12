package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/machinebox/progress"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func RunSSHCommand(sshClient *ssh.Client, command string) error {
	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("create ssh session failed %s", err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run("bash -l -c '" + command + "'")
	if err != nil {
		return fmt.Errorf("exec command %s failed %s", command, err)
	}
	return nil
}

func SftpLocalToRemote(sshClient *ssh.Client, localPath string, remotePath string) error {
	return Sftp(sshClient, localPath, remotePath, false)
}

func SftpRemoteToLocal(sshClient *ssh.Client, remotePath string, localPath string) error {
	return Sftp(sshClient, localPath, remotePath, true)
}

func Sftp(sshClient *ssh.Client, localPath string, remotePath string, remoteToLocal bool) error {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	if remoteToLocal {
		// 从远程复制到本地
		dstFile, err := os.Create(localPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// open source file
		srcFile, err := sftpClient.Open(remotePath)
		if err != nil {
			return err
		}
		fileStat, err := srcFile.Stat()
		if err != nil {
			return err
		}
		size := fileStat.Size()
		newSrcFile := progress.NewReader(srcFile)
		go func() {
			ctx := context.Background()
			progressChan := progress.NewTicker(ctx, newSrcFile, size, 5*time.Second)
			for p := range progressChan {
				fmt.Printf("copy to local remaining %s/%s %v ...\n", ByteCountIEC(p.N()), ByteCountIEC(size), p.Remaining().Round(time.Second))
			}
		}()

		bytes, err := io.Copy(dstFile, newSrcFile)
		if err != nil {
			return err
		}

		fmt.Printf("%s:%s => %s %d bytes copied\n", sshClient.RemoteAddr().String(), remotePath, localPath, bytes)
		// flush in-memory copy
		err = dstFile.Sync()
		if err != nil {
			return err
		}
	} else {
		//从本地复制到远程
		dstFile, err := sftpClient.Create(remotePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		srcFile, err := os.Open(localPath)
		if err != nil {
			return err
		}
		fileStat, err := srcFile.Stat()
		if err != nil {
			return err
		}
		size := fileStat.Size()
		newSrcFile := progress.NewReader(srcFile)
		go func() {
			ctx := context.Background()
			progressChan := progress.NewTicker(ctx, newSrcFile, size, 5*time.Second)
			for p := range progressChan {
				fmt.Printf("copy to remote remaining %s/%s %v ...\n", ByteCountIEC(p.N()), ByteCountIEC(size), p.Remaining().Round(time.Second))
			}
		}()

		bytes, err := io.Copy(dstFile, newSrcFile)
		if err != nil {
			return err
		}
		fmt.Printf("%s => %s:%s  %d bytes copied\n", localPath, sshClient.RemoteAddr().String(), remotePath, bytes)
	}
	return nil

}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

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

func ConnectSSh(config *ssh.ClientConfig, addr string) (*ssh.Client, error) {
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		//重试一次
		return ssh.Dial("tcp", addr, config)
	}
	return sshClient, err
}
