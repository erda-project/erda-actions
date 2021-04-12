package conf

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda/pkg/filehelper"
)

// detectLang 确保 conf.p 里的 Language/BuildType/ContainerType
func detectLang(p *params) (err error) {
	// modulePaths
	modulePaths := make([]string, 0, len(p.Modules))
	for _, m := range p.Modules {
		modulePaths = append(modulePaths, m.Path)
	}

	defer func() {
		if err != nil {
			err = errors.Errorf("failed to detect language/build_type/container_type, context: %s, modulePaths: %+v, err: %v",
				p.Context, modulePaths, err)
			return
		}
		bplog.Printf("language: %s, build_type: %s, container_type: %s\n",
			p.Language, p.BuildType, p.ContainerType)
	}()
	if p.FullLanguageInfo() {
		return nil
	}

	compatibleExplicitBpRepoVer(p)

	if p.FullLanguageInfo() {
		return nil
	}

	bplog.Printf("begin detect language/build_type/container_type\n")

	// check dir exist
	if err := filehelper.CheckExist(p.Context, true); err != nil {
		return errors.Errorf("context not found in fileSystem, err: %v", err)
	}
	err, result := langdetect.Detect(p.Context, modulePaths, types.DetectResult{
		Language:      p.Language,
		BuildType:     p.BuildType,
		ContainerType: p.ContainerType,
	})
	if err != nil {
		return err
	}
	p.Language = result.Language
	p.BuildType = result.BuildType
	p.ContainerType = result.ContainerType
	return nil
}
