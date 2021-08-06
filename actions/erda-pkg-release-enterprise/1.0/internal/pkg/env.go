package pkg

import (
	"path"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/pkg/errors"
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
		pkg.RepoToolsPath:   TmpRepoToolsPath(),
		pkg.DiceVersion:     config.ErdaVersion(),
	}

	// build erda install pkg before change erda to public or not
	patchVersion := path.Join(TmpRepoVersionPath(), config.ErdaVersion())
	exist, err := pkg.IsDirExists(patchVersion)
	if err != nil {
		return errors.WithMessage(err, "Judge patch version exists failed")
	}
	if exist {

		// erda release version dir exists in version repo before erda to public
		envs[pkg.ErdaToPublic] = pkg.False
	} else {

		// erda release version dir exists in oss archive bucket after erda to public
		envs[pkg.ErdaToPublic] = pkg.True
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

func TmpRepoToolsPath() string {
	return pkg.RepoTmpPath(config.RepoErdaTools())
}

func RepoToolsName() string {
	return pkg.RepoName(config.RepoErdaTools())
}
