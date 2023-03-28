package envconf

import (
	"os"
	"testing"
)

func TestGetTargetPlatforms(t *testing.T) {
	type arg struct {
		arch      string
		platforms string
	}
	testCases := []struct {
		name string
		arg  arg
		want string
	}{
		{
			name: "empty arch and platforms",
			want: "linux/amd64",
		},
		{
			name: "arm64 platform arch",
			arg: arg{
				arch: "arm64",
			},
			want: "linux/arm64",
		},
		{
			name: "empty arch but set platforms",
			arg: arg{
				platforms: "windows/amd64",
			},
			want: "windows/amd64",
		},
		{
			name: "both arch and platforms",
			arg: arg{
				arch:      "arm64",
				platforms: "linux/amd64,linux/arm64",
			},
			want: "linux/amd64,linux/arm64",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("DICE_ARCH", tc.arg.arch)
			os.Setenv("PLATFORMS", tc.arg.platforms)

			got := GetTargetPlatforms()
			if got != tc.want {
				t.Errorf("GetTargetPlatforms want: %s platforms, but got: %s", tc.want, got)
			}
		})
	}
}
