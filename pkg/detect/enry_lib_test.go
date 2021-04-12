package detect

import (
	"fmt"
	"testing"
)

func TestDetectDirLangs(t *testing.T) {
	var repo2 = "test-data/tomcattest"
	langs := DetectDirLangs(repo2)
	for _, lang := range langs {
		fmt.Println(lang.Type, lang.Percent)
	}
}

func TestDetectDirLangsByEnryLib(t *testing.T) {
	langs := detectDirLangsByEnryLib("/Users/sfwn/go/src/terminus.io/dice/dice/internal/pipeline/testdata/pampas-blog/services")
	for _, lang := range langs {
		fmt.Println(lang.Type, lang.Percent)
	}
}
