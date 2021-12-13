package run

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/build"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/build/buildartifact"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/pack"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
)

func Execute() error {

	var artifact *apistructs.BuildArtifact
	var packResult []byte

	defer func() {
		// store ${METAFILE}
		if artifact != nil {
			metaFileContent := buildartifact.MetadataFromArtifact(artifact)
			if err := filehelper.CreateFile(conf.PlatformEnvs().MetaFile, metaFileContent, 0644); err != nil {
				bplog.Println(err)
			}
		}
	}()

	// calc artifactSHA & identityText for query/register artifact
	artifactSHA, identityText, err := buildartifact.CalculateArtifactSHA()
	if err != nil {
		return errors.Errorf("failed to calculate SHA of artifact, err: %v", err)
	}

	artifact = handleForceBuildpack(artifactSHA, identityText)
	if artifact != nil {
		return nil
	}

	// prepare WorkDir for build & pack
	if err := prepareWorkDir(); err != nil {
		return errors.Errorf("failed to prepare workDir, err: %v", err)
	}

	// build
	if conf.Params().BuildkitEnable == "true" {
		if err := build.BuildkitBuild(); err != nil {
			return errors.Errorf("failed to build your code, err: %v", err)
		}
	} else {
		if err := build.Build(); err != nil {
			return errors.Errorf("failed to build your code, err: %v", err)
		}
	}

	if conf.Params().OnlyBuild {
		return nil
	}

	// pack
	if conf.Params().BuildkitEnable == "true" {
		packResult, err = pack.PackForBuildkit()
		if err != nil {
			return errors.Errorf("failed to pack image, err: %v", err)
		}
	} else {
		packResult, err = pack.Pack()
		if err != nil {
			return errors.Errorf("failed to pack image, err: %v", err)
		}
	}
	bplog.Printf("开始注册构建产物，SHA: %s\n", artifactSHA)
	newArtifact, err := buildartifact.RegisterBuildArtifact(artifactSHA, identityText, packResult)
	if err != nil {
		bplog.Printf("注册构建产物失败，请忽略。失败原因: %v\n", err)
	} else {
		bplog.Println("注册构建产物成功!")
	}
	artifact = newArtifact

	bplog.Println("打包成功!")

	return nil
}
