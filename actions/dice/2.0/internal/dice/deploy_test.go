package dice

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/erda-project/erda/pkg/envconf"
)

func TestInitEnv(t *testing.T) {
	var cfg conf
	err := envconf.Load(&cfg)
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg)
}

func getEnv() (*conf, error) {
	os.Setenv("DICE_APPLICATION_ID", "1")
	os.Setenv("DICE_CLUSTER_NAME", "terminus")
	os.Setenv("GITTAR_BRANCH", "feature/test")
	os.Setenv("DICE_OPERATOR_ID", "2")
	os.Setenv("DICE_WORKSPACE", "ENV")
	os.Setenv("DICE_ORG_ID", "3")
	os.Setenv("DICE_PROJECT_ID", "4")
	os.Setenv("DICE_OPENAPI_TOKEN", "xxxxxx")
	os.Setenv("DICE_OPENAPI_ADDR", "http://xxx.io")
	os.Setenv("DICE_DEPLOY_MODE", "RUNTIME")
	os.Setenv("PIPELINE_ID", "5")

	var cfg conf
	err := envconf.Load(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// func TestPrepareRequest(t *testing.T) {
// 	envCfg, err := getEnv()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(envCfg)

// 	pwd, _ := os.Getwd()

// 	f, err := os.Create(filepath.Join(pwd, "/dicehub_release"))
// 	defer os.Remove(filepath.Join(pwd, "/dicehub_release"))
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		f.Write([]byte("123456"))
// 	}

// 	req, err := prepareRequest(envCfg)
// 	t.Error(err)
// 	t.Log(req)
// }

func TestGetReleaseId(t *testing.T) {
	pwd, _ := os.Getwd()

	f, err := os.Create(filepath.Join(pwd, "/dicehub_release"))
	defer os.Remove(filepath.Join(pwd, "/dicehub_release"))
	if err != nil {
		t.Error(err.Error())
	} else {
		f.Write([]byte("123456"))
	}

	releaseId, _ := getReleaseId(pwd)
	t.Log(releaseId)
}

// func TestReportRuntimeId2PipelinePlatform(t *testing.T) {
// 	os.Setenv("DICE_OPENAPI_ADDR", "openapi.dev.terminus.io")
// 	os.Setenv("DICE_OPENAPI_TOKEN", "xxxx")
// 	os.Setenv("PIPELINE_ID", "169")
// 	os.Setenv("PIPELINE_TASK_ID", "694")
// 	require.NoError(t, reportRuntimeId2PipelinePlatform("10000", "20000"))
// }
