package utils

import (
	"strings"
	"fmt"
)

func ConvertType(t string) string {
	if t == "" {
		return t
	}
	return strings.ToUpper(fmt.Sprintf("%s_RELEASE", t))
}
