package ut

import (
	"math"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/pkg/detect"
	"github.com/erda-project/erda-actions/pkg/detect/bptype"
)

type Buildpack struct {
	// most likely language
	Language string `json:"language" yaml:"-"`
	BpRepo   string `json:"bp_repo" yaml:"bp_repo"`
	BpVer    string `json:"bp_ver" yaml:"bp_ver"`
}

func DetectBuildPack(detectPath string) (Buildpack, error) {
	langs := detect.DetectDirLangs(detectPath)
	if langs != nil && langs.Len() > 0 {
		inLangs, outLangs := splitLanguages(langs)
		for _, lang := range inLangs {
			support, bpRepo, bpVer := bptype.IsSupportedLanguage(lang.Type)
			if support {
				return Buildpack{
					Language: strings.ToLower(lang.Type),
					BpRepo:   bpRepo,
					BpVer:    bpVer,
				}, nil
			}
		}
		for _, lang := range outLangs {
			support, bpRepo, bpVer := bptype.IsSupportedLanguage(lang.Type)
			if support {
				return Buildpack{
					Language: strings.ToLower(lang.Type),
					BpRepo:   bpRepo,
					BpVer:    bpVer,
				}, nil
			}
		}
	} else {
		return Buildpack{}, errors.Errorf("cannot detect language detectPath: [%s]", detectPath)
	}

	return Buildpack{}, errors.New("Detect failed.")
}

func splitLanguages(langs detect.Languages) (detect.Languages, detect.Languages) {
	var inLangs detect.Languages
	var outLangs detect.Languages
	for _, lang := range langs {
		if lang.Internal {
			lang.Percent = math.Abs(lang.Percent)
			inLangs = append(inLangs, lang)
			sort.Sort(sort.Reverse(inLangs))
		} else {
			outLangs = append(outLangs, lang)
		}
	}
	return inLangs, outLangs
}
