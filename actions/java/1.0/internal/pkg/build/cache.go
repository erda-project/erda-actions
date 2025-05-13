package build

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/erda-project/erda-actions/actions/java/1.0/internal/pkg/conf"
	"github.com/erda-project/erda-actions/pkg/pack"
)

// handleCache 处理缓存
func handleCache(cfg conf.Conf) error {
	cacheKeyName := fmt.Sprintf("%s/%s/%s/%s", cfg.OrgName, cfg.ProjectName, cfg.AppName,
		strings.Replace(cfg.GittarBranch, "/", "-", -1))
	cacheStoragePath := cacheRootPath + "/" + cacheKeyName
	if !PathExists(cacheStoragePath) {
		err := os.MkdirAll(cacheStoragePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create dir %s err: %s", cacheStoragePath, err)
		}
	}
	cacheDirs := []string{"/root/.m2"}
	cacheMap := map[string]string{}
	for _, cacheDir := range cacheDirs {
		cacheFileName := base64.StdEncoding.EncodeToString([]byte(cacheDir)) + ".tar"
		cacheMap[cacheDir] = path.Join(cacheStoragePath, cacheFileName)
	}

	for cachePath, cacheTarFile := range cacheMap {
		// 之前运行存在缓存
		if PathExists(cacheTarFile) {
			os.MkdirAll(cachePath, os.ModePerm)
			fmt.Printf("start restore cache  %s\n", cacheTarFile)
			err := pack.UnTar(cacheTarFile, cachePath)
			if err != nil {
				return fmt.Errorf("failed to untar %s=>%s err: %s", cacheTarFile, cachePath, err)
			} else {
				fmt.Printf("restore cache  %s=>%s\n", cacheTarFile, cachePath)
			}
		} else {
			fmt.Printf("cacheFile:%s not exist\n", cacheTarFile)
		}
	}
	defer func() {
		for cacheDir, cacheTarFile := range cacheMap {
			if PathExists(cacheDir) {
				err := os.Chdir(cacheDir)
				fmt.Printf("pack cacheDir %s\n", cacheDir)
				if err != nil {
					fmt.Printf("failed to change cacheDir %s err: %s\n", cacheDir, err)
				}
				tmpFile := "/tmp/cache.tar"
				err = pack.Tar(tmpFile, ".")
				if err != nil {
					fmt.Printf("failed to tar %s=>%s err: %s\n", cacheDir, cacheTarFile, err)
				} else {
					fmt.Printf("success save cacheDir to %s", tmpFile)
				}
				mvCmd := exec.Command("mv", "-f", tmpFile, cacheTarFile)
				mvCmd.Stdout = os.Stdout
				mvCmd.Stderr = os.Stderr
				err = mvCmd.Run()
				if err != nil {
					fmt.Printf("failed to move %s=>%s err: %s\n", tmpFile, cacheTarFile, err)
				}
			} else {
				fmt.Printf("cacheDir %s not exist\n", cacheDir)
			}
		}
	}()

	return nil
}
