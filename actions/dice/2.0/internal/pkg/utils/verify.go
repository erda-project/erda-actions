package utils

import (
	"strings"

	"github.com/erda-project/erda/apistructs"
)

func VerifyWorkspace(workspace string) bool {
	switch strings.ToUpper(workspace) {
	case apistructs.WORKSPACE_DEV, apistructs.WORKSPACE_TEST,
		apistructs.WORKSPACE_STAGING, apistructs.WORKSPACE_PROD:
		return true
	default:
		return false
	}
}
