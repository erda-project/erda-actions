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
