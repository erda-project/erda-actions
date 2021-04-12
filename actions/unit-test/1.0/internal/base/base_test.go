package base

import (
	"testing"
)

func TestGetOsInfo(t *testing.T) {
	var (
		// uname -p or uname -s or uname -r
		osArch    = "p"
		osVersion = "r"
		osName    = "s"
	)
	t.Log(GetOsInfo(osArch))
	t.Log(GetOsInfo(osVersion))
	t.Log(GetOsInfo(osName))
}
