package pkg

import (
	"path"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-public/1.0/internal/config"
)

type Env struct {
}

func NewEnv() *Env {
	return &Env{}
}

// Set action env
func (e *Env) InitEnv() error {

	// env of repo and erda version
	var envs = map[string]string{
		pkg.RepoVersionPath: TmpRepoVersionPath(),
		pkg.RepoReleasePath: TmpRepoErdaReleasePath(),
		pkg.ErdaVersion:     config.ErdaVersion(),
	}

	// init git auth info, needed when building erda release pkg before erda to public
	envs[pkg.GitAccount] = config.GitInfo().Account
	envs[pkg.GitToken] = config.GitInfo().Token

	if err := pkg.InitEnv(envs); err != nil {
		return err
	}

	return nil
}

func TmpRepoVersionPath() string {
	return pkg.RepoTmpPath(config.RepoVersion())
}

func RepoVersionName() string {
	return pkg.RepoName(config.RepoVersion())
}

func TmpRepoErdaReleasePath() string {
	return pkg.RepoTmpPath(config.RepoErdaRelease())
}

func RepoErdaReleaseName() string {
	return pkg.RepoName(config.RepoErdaRelease())
}

func TmpErdaPatchPath() string {
	return pkg.GetPatchPath(config.ErdaVersion())
}

func GenErdaPublicName(osArch string) string {
	return pkg.GenErdaPublicName(config.ErdaVersion(), osArch)
}

func TmpOsArch(osArch string) string {
	return path.Join("/tmp", osArch)
}

func TmpOsArchRelease(osArch string) string {
	return path.Join("/tmp", osArch, RepoErdaReleaseName())
}
