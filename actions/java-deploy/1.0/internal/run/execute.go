package run

import (
	"os"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/dlog"
	"github.com/erda-project/erda/pkg/filehelper"
)

func Execute() error {
	// cd into deploy workdir
	if err := os.Chdir(conf.UserConf().Workdir); err != nil {
		return errors.Errorf("failed to cd into workdir, workdir: %s, err: %v\n", conf.UserConf().Workdir, err)
	}

	var (
		metaContent string
		err         error
	)
	if conf.UserConf().Cmd != "" {
		metaContent, err = executeCmd()
		if err != nil {
			return err
		}
	} else {
		metaContent, err = executeMvn()
		if err != nil {
			return err
		}
	}

	// metafile
	err = filehelper.CreateFile(conf.PlatformConf().MetaFile, metaContent, 0644)
	if err != nil {
		dlog.Printf("failed to create metadata file: %")
	}

	return nil
}
