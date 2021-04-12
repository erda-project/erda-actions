package langdetect

import (
	"reflect"
	"testing"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
)

func TestDetect(t *testing.T) {
	type args struct {
		dir         string
		modulePaths []string
		presets     []types.DetectResult
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 types.DetectResult
	}{
		{
			name: "java-multi-module",
			args: args{
				dir:         "../../../bp/java/build/maven/testdata/multi-modules",
				modulePaths: []string{"blog-service/blog-service-impl", "user-service/user-service-impl", "blog-web"},
				presets: []types.DetectResult{{
					Language:      "",
					BuildType:     "",
					ContainerType: "",
				}},
			},
			want: nil,
			want1: types.DetectResult{
				Language:      types.LanguageJava,
				BuildType:     types.BuildTypeMaven,
				ContainerType: types.ContainerTypeSpringBoot,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Detect(tt.args.dir, tt.args.modulePaths, tt.args.presets...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detect() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Detect() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
