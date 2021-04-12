package run

import (
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/pkg/detect/bptype"
	"github.com/erda-project/erda/pkg/filehelper"
)

// prepareWorkDir
//
// ${WORKDIR}
// - /bp/build
// - /bp/pack
// - /assets
// - /code
func prepareWorkDir() error {
	// copy correct build_type bp
	if err := filehelper.Copy(
		filepath.Join(conf.EasyUse().BpDir, string(conf.Params().Language), "build", string(conf.Params().BuildType)),
		conf.EasyUse().BpBuildTypeInWorkDir); err != nil {
		return errors.Errorf("failed to copy build_type bp, err: %v", err)
	}
	if err := bptype.RenderConfigToDir(conf.EasyUse().BpBuildTypeInWorkDir); err != nil {
		return errors.Errorf("failed to render build_type bp files: %v", err)
	}
	// copy correct container_type bp
	if err := filehelper.Copy(
		filepath.Join(conf.EasyUse().BpDir, string(conf.Params().Language), "pack", string(conf.Params().ContainerType)),
		conf.EasyUse().BpContainerTypeInWorkDir); err != nil {
		return errors.Errorf("failed to copy container_type bp, err: %v", err)
	}
	if err := bptype.RenderConfigToDir(conf.EasyUse().BpContainerTypeInWorkDir); err != nil {
		return errors.Errorf("failed to render container_type bp files: %v", err)
	}

	// copy assets
	if err := filehelper.Copy(conf.EasyUse().AssetsDir, conf.EasyUse().AssetsInWorkDir); err != nil {
		return errors.Errorf("failed to copy assets, err: %v", err)
	}

	// copy code
	if err := copy.Copy(conf.Params().Context, conf.EasyUse().CodeInWorkDir); err != nil {
		return err
	}

	//// copy ssh_config to wd.
	//if err := copy.Copy("/root/.ssh", filepath.Join(conf.PlatformEnvs().WorkDir, "ssh_config")); err != nil {
	//	return err
	//}

	return nil
}
