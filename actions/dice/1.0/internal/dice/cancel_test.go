package dice

import (
	"os"
	"testing"
)

func TestGetDiceInfo(t *testing.T) {
	err := storeDiceInfo("87649", "5852", "")
	defer os.Remove("diceinfo")
	if err != nil {
		t.Error(err)
	}

	deploymentId, runtimeId, err := getDiceInfo("")
	if err != nil {
		t.Error(err)
	}

	t.Log(deploymentId, runtimeId)
}

func TestCancelRequest(t *testing.T) {
	req := &CancelReq{
		RuntimeId: 5852,
	}

	os.Setenv("DICE_APPLICATION_ID", "1")
	os.Setenv("DICE_CLUSTER_NAME", "terminus")
	os.Setenv("GITTAR_BRANCH", "feature/test")
	os.Setenv("DICE_OPERATOR_ID", "2")
	os.Setenv("DICE_WORKSPACE", "ENV")
	os.Setenv("DICE_ORG_ID", "3")
	os.Setenv("DICE_PROJECT_ID", "4")
	os.Setenv("DICE_DEPLOY_MODE", "RUNTIME")
	os.Setenv("PIPELINE_ID", "5")
	os.Setenv("DICE_OPENAPI_TOKEN", "xxx")
	os.Setenv("DICE_OPENAPI_ADDR", "https://openapi.terminus.io")

	conf, err := getEnv()
	if err != nil {
		t.Error(err)
	}
	err = cancelRequest(req, "87649", conf)
	if err != nil {
		t.Error(err)
	}

	os.Clearenv()
}

func Test_Cancel(t *testing.T) {
	err := storeDiceInfo("87649", "5852", "")
	defer os.Remove("diceinfo")
	if err != nil {
		t.Error(err)
	}

	os.Setenv("DICE_APPLICATION_ID", "1")
	os.Setenv("DICE_CLUSTER_NAME", "terminus")
	os.Setenv("GITTAR_BRANCH", "feature/test")
	os.Setenv("DICE_OPERATOR_ID", "2")
	os.Setenv("DICE_WORKSPACE", "ENV")
	os.Setenv("DICE_ORG_ID", "3")
	os.Setenv("DICE_PROJECT_ID", "4")
	os.Setenv("DICE_DEPLOY_MODE", "RUNTIME")
	os.Setenv("PIPELINE_ID", "5")
	os.Setenv("DICE_OPENAPI_TOKEN", "xxxxx")
	os.Setenv("DICE_OPENAPI_ADDR", "https://openapi.terminus.io")

	err = Cancel()
	if err != nil {
		t.Error(err)
	}

	os.Clearenv()
}
