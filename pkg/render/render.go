package render

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

// RenderTemplate 渲染 dstDir 下模板文件，将{{xxx}}替换成给定内容
func RenderTemplate(dstDir string, cfgMap map[string]string) error {
	return filepath.Walk(dstDir, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Errorf("render template %s err, %v", path, err)
			}

			// 替换文件内容
			re := regexp.MustCompile("{{.+?}}")
			replacedContent := re.ReplaceAllStringFunc(string(bytes), func(s string) string {
				k := s[2:(len(s) - 2)]
				if v, ok := cfgMap[k]; ok {
					return v
				}
				return s
			})

			if replacedContent != string(bytes) { // 若文件含有{{}}占位符，则将替换后的内容回写文件
				if err := ioutil.WriteFile(path, []byte(replacedContent), info.Mode()); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
