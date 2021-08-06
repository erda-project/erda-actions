package pkg

import (
	"path"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/internal/bin"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/pkg/errors"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) PrepareRepo() error {

	// cp version to /tmp
	if err := pkg.CopyRepoToTmp(config.RepoVersion()); err != nil {
		return err
	}

	// cp tools to /tmp
	if err := pkg.CopyRepoToTmp(config.RepoErdaTools()); err != nil {
		return err
	}

	// tools execute script prepare
	toolsBuildScript := path.Join(TmpRepoToolsPath(), "build")
	if err := pkg.ReplaceFile(bin.PrivateExecuteScript, toolsBuildScript, 0666); err != nil {
		return errors.WithMessage(err, "replace build script in tools")
	}

	return nil
}
