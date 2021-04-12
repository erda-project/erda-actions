package md5util

import (
	"crypto/md5"
	"fmt"
	"path/filepath"

	"github.com/erda-project/erda/pkg/strutil"
)

// AppCache 生成应用缓存镜像名称
func AppCacheRepo(localRegistry, gitRepo, gitBranch, dstDir string) string {
	cacheID := fmt.Sprintf("%x", md5.Sum([]byte(strutil.Concat(gitRepo, gitBranch, dstDir))))

	return fmt.Sprintf("%s/appcache:latest", filepath.Join(localRegistry, cacheID))
}
