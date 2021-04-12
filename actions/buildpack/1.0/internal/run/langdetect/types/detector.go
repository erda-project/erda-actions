package types

import (
	"fmt"
)

const (
	LanguageJava       Language = "java"
	LanguageNode       Language = "node"
	LanguageDockerfile Language = "dockerfile"

	BuildTypeMaven          BuildType = "maven"
	BuildTypeMavenEdas      BuildType = "maven-edas"
	BuildTypeMavenEdasDubbo BuildType = "maven-edas-dubbo"
	BuildTypeNpm            BuildType = "npm"
	BuildTypeDockerfile     BuildType = "dockerfile"

	ContainerTypeSpringBoot     ContainerType = "springboot"
	ContainerTypeSpringBootAsia ContainerType = "springboot-asia"
	ContainerTypeEdas           ContainerType = "edas"
	ContainerTypeEdasDubbo      ContainerType = "edas-dubbo"
	ContainerTypeEdasDubboCnooc ContainerType = "edas-dubbo-cnooc"
	ContainerTypeHerd           ContainerType = "herd"
	ContainerTypeSpa            ContainerType = "spa"
	ContainerTypeDockerfile     ContainerType = "dockerfile"
)

type Language string
type BuildType string
type ContainerType string

type DetectResult struct {
	Language
	BuildType
	ContainerType
}

func (r DetectResult) OK() bool {
	return r.Language != "" && r.BuildType != "" && r.ContainerType != ""
}

type Detector interface {
	Language() Language
	BuildType() BuildType
	// param buildType may be auto detected or passed-in
	ContainerType(buildType BuildType) ContainerType

	SupportedBuildTypes() []BuildType
	SupportedContainerTypes() []ContainerType
}

func DescDetector(d Detector) string {
	return fmt.Sprintf("languge: %s, supportedBuildTypes: %v, supportedContainerTypes: %v",
		d.Language(), d.SupportedBuildTypes(), d.SupportedContainerTypes())
}
