package conf

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type easyUse struct {
	// 计算出来的缓存镜像
	CalculatedCacheImage string
	// 默认缓存镜像，当计算出来的缓存镜像下载失败时使用
	DefaultCacheImage string
	// Build 出来的镜像，从该镜像中拷贝出编译产物
	// 同时该镜像重新 tag 后会作为缓存镜像
	DockerImageFromBuild string
	// 代码相对地址
	RelativeCodeContext string

	OptActionDir             string
	AssetsDir                string
	AssetsInWorkDir          string
	AssetsJavaAgentDir       string
	AssetsJavaAgentInWorkDir string
	BpDir                    string
	BpBuildTypeInWorkDir     string
	BpContainerTypeInWorkDir string
	CodeInWorkDir            string
}

func generateEasyUseLast() (*easyUse, error) {
	var use easyUse

	use.OptActionDir = "/opt/action"
	use.AssetsDir = filepath.Join(use.OptActionDir, "assets")
	use.AssetsInWorkDir = filepath.Join(cfg.platformEnvs.WorkDir, "assets")
	use.AssetsJavaAgentDir = filepath.Join(use.AssetsDir, "java-agent")
	use.AssetsJavaAgentInWorkDir = filepath.Join(use.AssetsInWorkDir, "assets")
	use.BpDir = filepath.Join(use.OptActionDir, "bp")
	use.BpBuildTypeInWorkDir = filepath.Join(cfg.platformEnvs.WorkDir, "bp", "build")
	use.BpContainerTypeInWorkDir = filepath.Join(cfg.platformEnvs.WorkDir, "bp", "pack")
	use.CodeInWorkDir = filepath.Join(cfg.platformEnvs.WorkDir, "code")

	// RelativeCodeContext
	// ${git}/services -> git/services
	// 相对代码目录需要包含 git，因为 git/services 和 git2/services 是不同的
	relativeCodeContext, err := filepath.Rel(cfg.platformEnvs.ContextDir, cfg.params.Context)
	if err != nil {
		return nil, errors.Errorf("failed to get code relative path: %v", err)
	}
	use.RelativeCodeContext = relativeCodeContext

	cacheID := fmt.Sprintf("%x", md5.Sum([]byte(cfg.platformEnvs.ProjectAppAbbr+cfg.platformEnvs.GittarBranch+cfg.params.Context)))

	// CalculatedCacheImage
	use.CalculatedCacheImage = fmt.Sprintf("%s/cidepcache:latest",
		filepath.Join(cfg.platformEnvs.BpDockerCacheRegistry, cacheID))

	// DockerImageFromBuild
	use.DockerImageFromBuild = fmt.Sprintf("%s-%d", use.CalculatedCacheImage, time.Now().UnixNano())

	// DefaultCacheImage
	use.DefaultCacheImage = cfg.platformEnvs.DefaultDepCacheImage

	return &use, nil
}
