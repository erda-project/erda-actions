package pkg

import (
	"fmt"
	"os"
	"path"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-public/1.0/internal/bin"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-public/1.0/internal/config"
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

	// cp erda-release to /tmp
	if err := pkg.CopyRepoToTmp(config.RepoErdaRelease()); err != nil {
		return err
	}

	// erda-release execute script prepare
	buildPackageScript := path.Join(TmpRepoErdaReleasePath(), "build_package.sh")
	if err := pkg.ReplaceFile(bin.PublicExecuteScript, buildPackageScript, 0666); err != nil {
		return errors.WithMessage(err, "replace build_package script in erda-release")
	}

	// version to erda-release repo and prepare erda-release/version
	if err := r.ErdaReleaseVersionDeal(); err != nil {
		return errors.WithMessage(err, "prepare version to erda-release")
	}

	return nil
}

func (r *Repo) ErdaReleaseVersionDeal() error {

	// copy version repo to erda-release
	if _, err := pkg.ExecCmd(os.Stdout, os.Stderr, "", "cp", "-a",
		TmpRepoVersionPath(), TmpRepoErdaReleasePath()); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("copy %s to %s "+
			"repo failed", TmpRepoVersionPath(), TmpRepoErdaReleasePath()))
	}

	// copy erda patch version to version repo in erda
	erdaReleaseVersionPath := path.Join(TmpRepoErdaReleasePath(), RepoVersionName())
	if _, err := pkg.ExecCmd(os.Stdout, os.Stderr, "", "cp", "-a",
		TmpErdaPatchPath(), erdaReleaseVersionPath); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("copy %s to %s"+
			"repo failed", TmpErdaPatchPath(), erdaReleaseVersionPath))
	}

	// replace version compose script
	versionComposePath := path.Join(erdaReleaseVersionPath, "compose.sh")
	if err := pkg.ReplaceFile(bin.VersionComposeScript, versionComposePath, 0666); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("replace compose.sh in %s failed", versionComposePath))
	}

	return nil
}
