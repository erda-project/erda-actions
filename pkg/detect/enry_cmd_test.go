package detect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectDirLangsByEnryCmd(t *testing.T) {
	langs, err := detectDirLangsByEnryCmd("/Users/sfwn/go/src/git.terminus.io/spot/spot")
	assert.NoError(t, err)
	for _, lang := range langs {
		fmt.Println(lang.Type, lang.Percent)
	}
}
