package node

import (
	"path/filepath"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/dlog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	packageJson       = "package.json"
	packageLockJson   = "package-lock.json"
	yarnLock          = "yarn.lock"
	nginxConfTemplate = "nginx.conf.template"
)

type Node struct {
	dir string
}

func New(dir string) *Node {
	return &Node{dir: dir}
}

func (n Node) Language() types.Language {
	return types.LanguageNode
}

// npm
func (n Node) BuildType() types.BuildType {
	// npm
	// check package.json
	dlog.TryToFindFileUnderContextRoot(packageJson)
	if err := filehelper.CheckExist(filepath.Join(n.dir, packageJson), false); err == nil {
		dlog.FindFileUnderContextRoot(packageJson)
		return types.BuildTypeNpm
	}
	dlog.NotFoundFileUnderContextRoot(packageJson)

	return ""
}

// herd, spa
func (n Node) ContainerType(buildType types.BuildType) types.ContainerType {
	if buildType == types.BuildTypeNpm {

		// spa
		dlog.TryToFindFileUnderContextRoot(nginxConfTemplate)
		if err := filehelper.CheckExist(filepath.Join(n.dir, nginxConfTemplate), false); err == nil {
			dlog.FindFileUnderContextRoot(nginxConfTemplate)
			return types.ContainerTypeSpa
		}
		dlog.NotFoundFileUnderContextRoot(nginxConfTemplate)

		// default is herd
		bplog.Printf("use default container_type: %s", types.ContainerTypeHerd)
		return types.ContainerTypeHerd
	}

	return ""
}

func (n Node) SupportedBuildTypes() []types.BuildType {
	return []types.BuildType{types.BuildTypeNpm}
}

func (n Node) SupportedContainerTypes() []types.ContainerType {
	return []types.ContainerType{types.ContainerTypeHerd, types.ContainerTypeSpa}
}
