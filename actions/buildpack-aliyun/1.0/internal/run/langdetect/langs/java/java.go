package java

import (
	"io/ioutil"
	"path/filepath"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/dlog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
)

const (
	pomXML     = "pom.xml"
	springBoot = "spring-boot"
)

type Java struct {
	dir         string
	modulePaths []string
}

func New(dir string, modulePaths ...string) *Java {
	return &Java{dir: dir, modulePaths: modulePaths}
}

func (j Java) Language() types.Language {
	return types.LanguageJava
}

// maven, maven-edas
func (j Java) BuildType() types.BuildType {
	// maven
	dlog.TryToDetectExpectedBuildType(types.BuildTypeMaven)
	// check parent pom.xml
	dlog.TryToFindFileUnderContextRoot(pomXML)
	pomXMLFilePath := filepath.Join(j.dir, pomXML)
	if err := filehelper.CheckExist(pomXMLFilePath, false); err == nil {
		dlog.FindFileInPath(pomXML, pomXMLFilePath)
		return types.BuildTypeMaven
	}
	dlog.NotFoundFileUnderContextRoot(pomXML)

	// TODO gradle

	return ""
}

// springboot, edas
func (j Java) ContainerType(buildType types.BuildType) types.ContainerType {
	// maven: springboot
	if buildType == types.BuildTypeMaven {
		var allPomPaths []string
		// parent pom
		allPomPaths = append(allPomPaths, filepath.Join(j.dir, pomXML))
		// module poms
		for _, modulePath := range j.modulePaths {
			allPomPaths = append(allPomPaths, filepath.Join(j.dir, modulePath, pomXML))
		}
		// judge
		for _, pomPath := range allPomPaths {
			dlog.TryToFindContentInFile(springBoot, pomPath)
			if isSpringBoot(pomPath) {
				dlog.FindContentInFile(springBoot, pomPath)
				return types.ContainerTypeSpringBoot
			}
			dlog.NotFoundContentInFile(springBoot, pomPath)
		}
	}

	// maven-edas: edas
	if buildType == types.BuildTypeMavenEdas {
		bplog.Printf("build_type is %q, so container_type is %q\n", types.BuildTypeMavenEdas, types.ContainerTypeEdas)
		return types.ContainerTypeEdas
	}

	return ""
}

func (j Java) SupportedBuildTypes() []types.BuildType {
	return []types.BuildType{types.BuildTypeMaven, types.BuildTypeMavenEdas, types.BuildTypeMavenEdasDubbo}
}

func (j Java) SupportedContainerTypes() []types.ContainerType {
	return []types.ContainerType{types.ContainerTypeSpringBoot, types.ContainerTypeEdas, types.ContainerTypeEdasDubbo, types.ContainerTypeEdasDubboCnooc, types.ContainerTypeSpringBootAsia}
}

func isSpringBoot(pomPath string) bool {
	pomXMLContent, err := ioutil.ReadFile(filepath.Join(pomPath))
	if err == nil && strutil.Contains(string(pomXMLContent), springBoot) {
		return true
	}
	return false
}
