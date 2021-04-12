package langdetect

import (
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/langs/dockerfile"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/langs/java"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/langs/node"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
)

// supportedBuildTypes 保存所有 buildType，格式为 map[Language]map[BuildType]struct{}
var supportedBuildTypes map[types.Language]map[types.BuildType]struct{}

// supportedContainerTypes 保存所有 containerType，格式为 map[Language]map[ContainerType]struct{}
var supportedContainerTypes map[types.Language]map[types.ContainerType]struct{}

func init() {
	supportedBuildTypes = make(map[types.Language]map[types.BuildType]struct{})
	supportedContainerTypes = make(map[types.Language]map[types.ContainerType]struct{})

	allDetectors := []types.Detector{
		java.New(""), node.New(""), dockerfile.New(""),
	}
	for _, d := range allDetectors {
		// build type
		for _, buildType := range d.SupportedBuildTypes() {
			if supportedBuildTypes[d.Language()] == nil {
				supportedBuildTypes[d.Language()] = make(map[types.BuildType]struct{})
			}
			supportedBuildTypes[d.Language()][buildType] = struct{}{}
		}

		// container type
		for _, containerType := range d.SupportedContainerTypes() {
			if supportedContainerTypes[d.Language()] == nil {
				supportedContainerTypes[d.Language()] = make(map[types.ContainerType]struct{})
			}
			supportedContainerTypes[d.Language()][containerType] = struct{}{}
		}
	}
}
