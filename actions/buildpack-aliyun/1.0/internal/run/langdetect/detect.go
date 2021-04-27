package langdetect

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/langs/dockerfile"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/langs/java"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/langs/node"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/types"
)

// Detect 根据 dir 自动检测 language/build_type/container_type.
// @presets: 用户指定的类型
func Detect(dir string, modulePaths []string, presets ...types.DetectResult) (error, types.DetectResult) {

	allDetectors := []types.Detector{
		java.New(dir, modulePaths...), node.New(dir), dockerfile.New(dir),
	}

	var preset *types.DetectResult
	if len(presets) > 0 && presets[0].Language != "" {
		preset = &presets[0]

		// check preset
		if preset.BuildType != "" {
			_, ok := supportedBuildTypes[preset.Language][preset.BuildType]
			if !ok {
				return errors.Errorf("not match! language: %s, build_type: %s", preset.Language, preset.BuildType), types.DetectResult{}
			}
		}
		if preset.ContainerType != "" {
			_, ok := supportedContainerTypes[preset.Language][preset.ContainerType]
			if !ok {
				return errors.Errorf("not match! language: %s, container_type: %s", preset.Language, preset.BuildType), types.DetectResult{}
			}
		}
	}

	for _, detector := range allDetectors {
		matched, result := matchOneDetector(detector, preset)
		if matched {
			return nil, result
		}
	}

	return errors.Errorf("no detector matched!"), types.DetectResult{}
}

// matchOneDetector return one detector's detect result with preset DetectResult.
func matchOneDetector(d types.Detector, preset *types.DetectResult) (bool, types.DetectResult) {
	// check language
	if preset != nil {
		if preset.Language != d.Language() {
			return false, types.DetectResult{}
		}
	}

	bplog.Printf("try detector, %s\n", types.DescDetector(d))

	var result types.DetectResult

	// language
	result.Language = d.Language()

	// buildType
	result.BuildType = d.BuildType()
	if preset != nil && preset.BuildType != "" {
		result.BuildType = preset.BuildType
	}
	if result.BuildType == "" {
		bplog.Printf("no suitable build_type detected, skip detect with %s detector\n", d.Language())
		return false, result
	}
	bplog.Printf("detected build_type: %s\n", result.BuildType)

	// containerType
	result.ContainerType = d.ContainerType(result.BuildType)
	if preset != nil && preset.ContainerType != "" {
		result.ContainerType = preset.ContainerType
	}
	if result.ContainerType == "" {
		bplog.Printf("no suitable container_type detected, skip detect with %s detector\n", d.Language())
		return false, result
	}
	bplog.Printf("detected container_type: %s\n", result.ContainerType)

	return result.OK(), result
}
