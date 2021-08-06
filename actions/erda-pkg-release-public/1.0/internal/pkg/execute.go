package pkg

import (
	"fmt"
	"os"
	"path"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-public/1.0/internal/config"
	"github.com/pkg/errors"
)

var osArches = []string{
	"linux-x86",
}

func Execute() error {

	// oss init
	oss := pkg.NewOSS(config.OssInfo(), config.ErdaVersion(), config.ReleaseType(),
		pkg.OssPkgReleasePublicPath, false)
	if err := oss.InitOssConfig(); err != nil {
		return err
	}

	// prepare erda patch version info to /tmp/
	if err := oss.PreparePatchRelease(); err != nil {
		return err
	}

	// prepare repo to use
	repo := NewRepo()
	if err := repo.PrepareRepo(); err != nil {
		return err
	}

	// set init env
	var env = NewEnv()
	if err := env.InitEnv(); err != nil {
		return err
	}

	// tool-pack execute
	releasePkgPathInfo, releasePkgInfo, err := ErdaPkgRelease()
	if err != nil {
		return err
	}

	// upload release install pkg of erda
	if err := oss.ReleasePackage(releasePkgPathInfo); err != nil {
		return err
	}

	// write metafile
	metafile := pkg.NewMetafile(oss, config.MetaFile())
	if err := metafile.WriteMetaFile(releasePkgInfo); err != nil {
		return err
	}

	return nil
}

// ErdaPkgRelease to build some erda installing package with
// some version specified by ERDA_VERSION
func ErdaPkgRelease() (map[string]string, map[string]string, error) {

	var (
		err            error
		releasePkgInfo map[string]string
	)

	// build erda release package
	if releasePkgInfo, err = PublicPkgRelease(); err != nil {
		return nil, nil, err
	}

	// erda release package path info
	releasePkgPathInfo := map[string]string{}
	for osArch, erdaPkg := range releasePkgInfo {
		releasePkgPathInfo[osArch] = path.Join(TmpOsArchRelease(osArch), "package", erdaPkg)
	}

	return releasePkgPathInfo, releasePkgInfo, nil
}

func PublicPkgRelease() (map[string]string, error) {

	releasePkgInfo := map[string]string{}

	// build public release package of erda according to orArch
	for _, osArch := range osArches {

		// prepare osArch erda-release
		tmpOsArch := TmpOsArch(osArch)
		tmpOsArchRelease := TmpOsArchRelease(osArch)

		if _, err := pkg.ExecCmd(os.Stdout, os.Stderr, "", "mkdir", "-p", tmpOsArch); err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("create tmp dir /tmp/%s", osArch))
		}

		if _, err := pkg.ExecCmd(os.Stdout, os.Stderr, "", "cp", "-rf",
			TmpRepoErdaReleasePath(), tmpOsArch); err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("cp release pkg to /tmp/%s", osArch))
		}

		// replace build script
		if _, err := pkg.ExecCmd(os.Stdout, os.Stderr, tmpOsArchRelease, "bash", "-x", "build_package.sh"); err != nil {
			return nil, errors.WithMessage(err, "build public erda install package")
		}

		// archive pkg release package of erda specified by osArch
		pkgRelatedPath := fmt.Sprintf("package/%s", GenErdaPublicName(osArch))
		if _, err := pkg.ExecCmd(os.Stdout, os.Stderr, tmpOsArchRelease, "tar", "-cvzf",
			pkgRelatedPath, "erda/"); err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("archive erda release package "+
				"%s related to %s", GenErdaPublicName(osArch), osArch))
		}

		releasePkgInfo[osArch] = GenErdaPublicName(osArch)
	}

	return releasePkgInfo, nil
}
