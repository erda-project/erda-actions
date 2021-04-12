package dockerfile

import (
	"path/filepath"

	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/dlog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
)

const (
	FileDockerfile = "Dockerfile"
)

type Dockerfile struct {
	dir            string
	dockerfilePath string
}

func New(dir string, ops ...Option) *Dockerfile {
	d := Dockerfile{dir: dir, dockerfilePath: FileDockerfile}

	for _, op := range ops {
		op(&d)
	}

	return &d
}

type Option func(*Dockerfile)

func WithDockerfilePath(customPath string) Option {
	return func(d *Dockerfile) {
		d.dockerfilePath = customPath
	}
}

func (d Dockerfile) Language() types.Language {
	return types.LanguageDockerfile
}

// dockerfile
func (d Dockerfile) BuildType() types.BuildType {
	// dockerfile
	dlog.TryToFindFileInPath(FileDockerfile, d.dockerfilePath)
	if err := filehelper.CheckExist(filepath.Join(d.dir, d.dockerfilePath), false); err == nil {
		dlog.FindFileInPath(FileDockerfile, d.dockerfilePath)
		return types.BuildTypeDockerfile
	}
	dlog.NotFoundFileInPath(FileDockerfile, d.dockerfilePath)

	return ""
}

// dockerfile
func (d Dockerfile) ContainerType(buildType types.BuildType) types.ContainerType {
	if buildType == types.BuildTypeDockerfile {
		bplog.Printf("build_type is %q, so container_type is %q\n", types.BuildTypeDockerfile, types.ContainerTypeDockerfile)
		return types.ContainerTypeDockerfile
	}

	return ""
}

func (d Dockerfile) SupportedBuildTypes() []types.BuildType {
	return []types.BuildType{types.BuildTypeDockerfile}
}

func (d Dockerfile) SupportedContainerTypes() []types.ContainerType {
	return []types.ContainerType{types.ContainerTypeDockerfile}
}
