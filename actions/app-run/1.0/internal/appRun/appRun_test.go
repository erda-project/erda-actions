package appRun

import (
	"testing"
)

func TestHandlerPipelineYmlName(t *testing.T) {
	var tables = []struct {
		ymlName string
		result  string
	}{
		{
			"pipeline.yml",
			"pipeline.yml",
		},
		{
			"aa.yml",
			".dice/pipelines/aa.yml",
		},
		{
			".dice/pipelines/bb.yml",
			".dice/pipelines/bb.yml",
		},
		{
			".dice/bb.yml",
			".dice/pipelines/.dice/bb.yml",
		},
	}

	for _, data := range tables {
		var result = handlerPipelineYmlName(data.ymlName)
		if result != data.result {
			t.Fatalf("result not same like")
		}
	}
}

func Test_handlerPipelineYmlName(t *testing.T) {
	type args struct {
		ymlName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test .dice",
			args: args{
				ymlName: ".dice/pipelines/aa.yml",
			},
			want: ".dice/pipelines/aa.yml",
		},
		{
			name: "test .erda",
			args: args{
				ymlName: ".erda/pipelines/aa.yml",
			},
			want: ".erda/pipelines/aa.yml",
		},
		{
			name: "test default .dice",
			args: args{
				ymlName: "aa.yml",
			},
			want: ".dice/pipelines/aa.yml",
		},
		{
			name: "test pipeline.yml",
			args: args{
				ymlName: "pipeline.yml",
			},
			want: "pipeline.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlerPipelineYmlName(tt.args.ymlName); got != tt.want {
				t.Errorf("handlerPipelineYmlName() = %v, want %v", got, tt.want)
			}
		})
	}
}
