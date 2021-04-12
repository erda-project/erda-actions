package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/envconf"
)

type DeleteNodesRequest struct {
	AccessKey      string `env:"ACTION_AK" required:"true"`
	SecretKey      string `env:"ACTION_SK" required:"true"`
	Region         string `env:"ACTION_REGION" required:"true"`
	ScalingGroupId string `env:"ACTION_SCALING_GROUP_ID" required:"true"`

	InstanceIDs string `env:"ACTION_INSTANCE_IDS"`

	IsCron          bool   `env:"ACTION_IS_CRON"`
	InstanceIDsFile string `env:"ACTION_INSTANCE_IDS_FILE"`
}

func main() {
	err := DeleteNodes()
	if err != nil {
		os.Exit(1)
	}
}

func DeleteNodes() error {
	req := DeleteNodesRequest{}

	if err := envconf.Load(&req); err != nil {
		err := errors.Errorf("failed to parse platform envs: %v", err)
		return err
	}

	client, err := ess.NewClientWithAccessKey(req.Region, req.AccessKey, req.SecretKey)
	if err != nil {
		err := errors.Errorf("create client error: %v", err)
		return err
	}
	request := ess.CreateRemoveInstancesRequest()
	request.Scheme = "https"

	// if is cron job, get instance ids from file passed from previous action
	if req.IsCron {
		_, err := os.Stat(req.InstanceIDsFile)
		if err != nil {
			logrus.Errorf("cron job, but do not have meta file: %s, error: %v", req.InstanceIDsFile, err)
		}
		content, err := ioutil.ReadFile(req.InstanceIDsFile)
		if err != nil {
			logrus.Errorf("cron job, read file failed: %s, error: %v", req.InstanceIDsFile, err)
		}
		req.InstanceIDs = string(content)
	}

	if isEmpty(req.ScalingGroupId) {
		err := errors.Errorf("invalid request: %v, error: %v", req, err)
		return err
	}

	if isEmpty(req.InstanceIDs) {
		logrus.Infof("empty scaling group")
		return nil
	}

	request.ScalingGroupId = req.ScalingGroupId
	instanceIds := strings.Split(req.InstanceIDs, ",")
	request.InstanceId = &instanceIds
	request.RemovePolicy = "release"

	response, err := client.RemoveInstances(request)
	if err != nil {
		err := errors.Errorf("failed to delete instances: %v, error: %v", request, err)
		return err
	}
	fmt.Printf("delete instances success, response: %v", response)
	return nil
}

func isEmpty(str string) bool {
	return strings.Replace(str, " ", "", -1) == ""
}
