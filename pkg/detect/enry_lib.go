package detect

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/src-d/enry.v1"

	"github.com/erda-project/erda-actions/pkg/detect/bptype"
)

// Language constructed by Language type, percent, internal
type Language struct {
	Type     string  `json:"type"`
	Percent  float64 `json:"percent"`
	Internal bool    `json:"internal,omitempty"`
}

type Languages []Language

func (langs Languages) Len() int { return len(langs) }

func (langs Languages) Less(i, j int) bool { return langs[i].Percent < langs[j].Percent }

func (langs Languages) Swap(i, j int) { langs[i], langs[j] = langs[j], langs[i] }

// DetectDirLangs detect using enry cmd first;
// if err, then use enry lib.
// NOTE: absolute path is expected.
// UPDATE: using enry lib directly.
func DetectDirLangs(dir string) Languages {
	//langs, err := detectDirLangsByEnryCmd(dir)
	//if err == nil {
	//	return langs
	//}
	return detectDirLangsByEnryLib(dir)
}

// detectDirLangsByEnryLib using enry go lib,
// not very accurate compare to enry cmd.
// TODO bug sql detect
func detectDirLangsByEnryLib(dir string) Languages {

	langMap := make(map[string]float64, 0)

	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if !f.Mode().IsDir() && !f.Mode().IsRegular() {
			return nil
		}

		if isDiceLogo := func(path string, f os.FileInfo) bool {
			if f.IsDir() {
			} else {
				// SPA
				if f.Name() == bptype.DICE_SPA_MARK {
					langMap[bptype.DICE_SPA] += -float64(f.Size())
					return true
				}
				// Herd
				if f.Name() == bptype.HERD_MARK {
					content, err := readFile(path, -1)
					if err != nil {
						return false
					}
					if strings.Contains(string(content), "herd ") {
						langMap[bptype.HERD] += -float64(f.Size())
						return true
					}
				}
				// root of dir
				if filepath.Dir(path) == dir {
					if f.Name() == bptype.DICE_DOCKERFILE_MARK {
						content, err := readFile(path, -1)
						if err != nil {
							return false
						}
						if strings.Contains(string(content), "#!dice") ||
							strings.Contains(string(content), "# dice-tags: dice-buildpack-dockerfile") {
							langMap[bptype.DICE_DOCKERFILE] += -float64(f.Size())
							return true
						}
					}
					if f.Name() == bptype.TOMCAT_MARK {
						content, err := readFile(path, -1)
						if err != nil {
							return false
						}
						if strings.Contains(string(content), "<packaging>war</packaging>") {
							langMap[bptype.TOMCAT] += -float64(f.Size())
							return true
						}
					}
				}
			}
			return false
		}(path, f); isDiceLogo {
			return nil
		}

		relativePath, err := filepath.Rel(dir, path)
		if err != nil {
			return nil
		}

		if relativePath == "." {
			return nil
		}

		if f.IsDir() {
			relativePath = relativePath + "/"
		}

		if enry.IsVendor(relativePath) || enry.IsDotFile(relativePath) ||
			enry.IsDocumentation(relativePath) || enry.IsConfiguration(relativePath) {
			if f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if f.IsDir() {
			return nil
		}

		lang, ok := enry.GetLanguageByExtension(path)
		if !ok {
			if lang, ok = enry.GetLanguageByFilename(path); !ok {
				content, err := readFile(path, int64(16*1024*1024))
				if err != nil {
					return nil
				}

				lang := enry.GetLanguage(filepath.Base(path), content)
				if lang == enry.OtherLanguage {
					return nil
				}
			}
		}

		// NOTE: check enry.IsConfiguration etc above
		canForcePass := func(lang string) bool {
			switch strings.ToLower(lang) {
			case "dockerfile":
				return true
			default:
				return false
			}
		}

		// only programming Language
		if enry.GetLanguageType(lang) != enry.Programming {
			// some Language is not programming Language but pass
			if !canForcePass(lang) {
				return nil
			}
		}

		langMap[lang] += float64(f.Size())

		return nil
	})

	// get total count
	var total float64
	for _, count := range langMap {
		if count > 0 {
			total += count
		}
	}

	var langs Languages
	for lang, count := range langMap {
		per, err := strconv.ParseFloat(fmt.Sprintf("%.4f", count/total*100.00), 64)
		if err != nil {
			per = float64(count) / float64(total) * 100.00
		}
		langs = append(langs, Language{Type: lang, Percent: per, Internal: bptype.IsInternalLang(lang)})
	}

	// order by Language.Percent desc
	sort.Sort(sort.Reverse(langs))

	return langs
}

func readFile(path string, limit int64) ([]byte, error) {
	if limit <= 0 {
		return ioutil.ReadFile(path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	if limit > 0 && size > limit {
		size = limit
	}
	buf := bytes.NewBuffer(nil)
	buf.Grow(int(size))
	_, err = io.Copy(buf, io.LimitReader(f, limit))
	return buf.Bytes(), err
}
