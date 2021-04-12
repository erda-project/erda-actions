package bptype

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/erda-project/erda-actions/pkg/config"
	"github.com/erda-project/erda/pkg/filehelper"

	"github.com/sirupsen/logrus"
)

// ----------
// internal
// ----------

// DICE_SPA
var (
	DICE_SPA         = "DICE_SPA"
	DICE_SPA_MARK    = "nginx.conf.template"
	DICE_SPA_BP_REPO = ""
)

// Herd
var (
	HERD         = "Herd"
	HERD_MARK    = "package.json"
	HERD_BP_REPO = ""
)

// DICE_DOCKERFILE
var (
	DICE_DOCKERFILE         = "Dice_Dockerfile"
	DICE_DOCKERFILE_MARK    = "Dockerfile"
	DICE_DOCKERFILE_BP_REPO = ""
)

var (
	TOMCAT             = "tomcat"
	TOMCAT_MARK        = "pom.xml"
	TOMCAT_BP_REPO     = ""
	TOMCAT_DEFAULT_VER = "feature/tomcat"
)

// ----------
// general supported
// ----------

// java
var (
	JAVA_BP_REPO = ""
)

// golang
var (
	GO_BP_REPO = ""
)

var configMap = map[string]string{}
var cfg config.PrivateDeployConfig

func init() {
	cfg = config.GetPrivateDeployConfig()
	repoPrefix := "file:///opt/action/buildpacks"
	config.GetPrivateDeployConfigMap()
	DICE_SPA_BP_REPO = repoPrefix + "/dice-bpack-termspa.git"
	HERD_BP_REPO = repoPrefix + "/dice-bpack-termnodejs.git"
	DICE_DOCKERFILE_BP_REPO = repoPrefix + "/dice-bpack-dockerimage.git"
	JAVA_BP_REPO = repoPrefix + "/dice-bpack-termjava.git"
	GO_BP_REPO = repoPrefix + "/dice-bpack-termgolang.git"
	TOMCAT_BP_REPO = repoPrefix + "/dice-bpack-termjava.git"

	configMap = config.GetPrivateDeployConfigMap()
}

// IsInternalLang check a language is internal or not
var IsInternalLang = func(lang string) bool {
	switch lang {
	case DICE_SPA, HERD, DICE_DOCKERFILE, TOMCAT:
		return true
	}
	return false
}

func IsSupportedLanguage(lang string) (supported bool, bpRepo, bpVer string) {

	polishLang := func(lang string) string {
		switch strings.ToLower(lang) {
		case "go":
			return "golang"
		case strings.ToLower(HERD):
			return "nodejs"
		case strings.ToLower("kotlin"):
			return "java"
		case strings.ToLower("dockerfile"):
			return "dockerimage"
		default:
			return strings.ToLower(lang)
		}
	}

	bpRepoPrefix := "file:///opt/action/buildpacks"
	bpDefaultVer := "master"

	checkBpOnTerminusGit := func(lang string) (supported bool, bpRepo, bpVer string) {
		bprepo_prefix := bpRepoPrefix + "/dice-bpack-"
		if lang != "dockerimage" {
			bprepo_prefix = bprepo_prefix + "term"
		}
		bpRepo = fmt.Sprintf("%s%s", bprepo_prefix, polishLang(lang))
		return true, fmt.Sprintf("%s.git", bpRepo), bpDefaultVer
	}

	switch l := polishLang(lang); l {
	case strings.ToLower(DICE_SPA):
		return true, DICE_SPA_BP_REPO, bpDefaultVer
	case strings.ToLower(HERD):
		return true, HERD_BP_REPO, bpDefaultVer
	case strings.ToLower(DICE_DOCKERFILE):
		return true, DICE_DOCKERFILE_BP_REPO, bpDefaultVer
	case "java":
		return true, JAVA_BP_REPO, bpDefaultVer
	case "golang":
		return true, GO_BP_REPO, bpDefaultVer
	case "tomcat":
		return true, TOMCAT_BP_REPO, TOMCAT_DEFAULT_VER

	default:
		return checkBpOnTerminusGit(l)
	}
}

func RenderConfigToDir(projectDir string) error {
	// 只遍历一层
	files, err := ioutil.ReadDir(projectDir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		if !fi.IsDir() {
			filePath := filepath.Join(projectDir, fi.Name())
			bytes, err := ioutil.ReadFile(filePath)
			if err == nil {
				result, change := RenderConfig(string(bytes))
				if change {
					logrus.Infof("render template file success, name: %s, path: %s", fi.Name(), filePath)
					err := filehelper.CreateFile(filePath, result, fi.Mode())
					if err != nil {
						logrus.Errorf("save template file success, name: %s, path: %s", fi.Name(), filePath)
					}
				}
			} else {
				logrus.Errorf("read file error: %v", err)
			}
		}
	}
	return nil
}

func RenderConfig(template string) (string, bool) {
	compile, _ := regexp.Compile("{{.+?}}")
	hasChange := false
	result := compile.ReplaceAllStringFunc(template, func(s string) string {
		key := s[2:(len(s) - 2)]
		value, ok := configMap[key]
		if ok {
			hasChange = true
			return value
		} else {
			v := os.Getenv(key)
			if v != "" {
				hasChange = true
				return v
			}
		}
		return s
	})
	return result, hasChange
}
