package main

import (
	"os"
	"testing"
)

func TestDeleteNodes(t *testing.T) {
	os.Setenv("ACTION_AK", "xxx")
	os.Setenv("ACTION_SK", "xxx")
	os.Setenv("ACTION_REGION", "cn-hangzhou")
	os.Setenv("ACTION_SCALING_GROUP_ID", "asg-bp1i65jd06uyyxym3b7v")
	os.Setenv("ACTION_INSTANCE_IDS", "")
	os.Setenv("ACTION_IS_CRON", "true")
	os.Setenv("ACTION_INSTANCE_IDS_FILE", "/tmp/instance_ids")
	err := DeleteNodes()
	if err != nil {
		panic(err)
	}
}
