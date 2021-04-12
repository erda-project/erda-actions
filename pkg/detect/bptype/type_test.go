package bptype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedLanguage(t *testing.T) {
	lang := "javascript"
	support, bpRepo, bpVer := IsSupportedLanguage(lang)
	require.False(t, support)
	fmt.Println(support)
	fmt.Println(bpRepo)
	fmt.Println(bpVer)
}
