package pkg

import (
	"fmt"
	"os"
	"path"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/pkg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var osArches = []string{
	"linux-x86",
}

func Execute() error {

	// oss init
	oss := pkg.NewOSS(config.OssInfo(), config.ErdaVersion(), config.ReleaseType(),
		pkg.OssPkgReleasePrivatePath, false)
	if err := oss.InitOssConfig(); err != nil {
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

	// prepare patch release of some version specified by erda version
	// when erda version before erda be public
	if os.Getenv(pkg.ErdaToPublic) == pkg.True {
		logrus.Infof("erda has push to public of this version %s. "+
			"Prepare patch version info", config.ErdaVersion())
		if err := oss.PreparePatchRelease(); err != nil {
			return err
		}
	} else {
		logrus.Infof("erda has not push to public of this version %s. "+
			"no need to prepare patch version info", config.ErdaVersion())
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

	if err := EnterprisePkgRelease(); err != nil {
		return nil, nil, err
	}

	releasePkgPathInfo := map[string]string{}
	releasePkgInfo := map[string]string{}
	for _, osArch := range osArches {

		// enterprise
		releasePkgPathInfo[osArch] = path.Join(TmpRepoToolsPath(), fmt.Sprintf(
			"dice-tools.%s.tar.gz", config.ErdaVersion()))

		releasePkgInfo[osArch] = fmt.Sprintf("dice-tools.%s.tar.gz", config.ErdaVersion())
	}

	return releasePkgPathInfo, releasePkgInfo, nil
}

func EnterprisePkgRelease() error {

	// replace build script
	_, err := pkg.ExecCmd(os.Stdout, os.Stderr, TmpRepoToolsPath(), "bash", "-x", "build", "pack")
	if err != nil {
		return errors.WithMessage(err, "build enterprise erda install package")
	}

	// doing something with offline release
	if config.ReleaseType() == pkg.ReleaseOffline {
		logrus.Info("release erda package for offline cluster...")
		_, err = pkg.ExecCmd(os.Stdout, os.Stderr, TmpRepoToolsPath(), "bash", "-x", "offline/build.sh")
		if err != nil {
			return errors.WithMessage(err, "deal with offline build script")
		}
		logrus.Info("release erda package for offline cluster successfully!!!")
	}

	// compress enterprise install package of erda
	_, err = pkg.ExecCmd(os.Stdout, os.Stderr, TmpRepoToolsPath(), "make", "tar")
	if err != nil {
		return errors.WithMessage(err, "build enterprise erda install package")
	}

	return nil
}
