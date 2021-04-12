package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	api "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/sirupsen/logrus"
)

type Ess struct {
	Region         string `env:"ACTION_REGION" required:"true"`
	AccessKeyID    string `env:"ACTION_AK" required:"true"`
	AccessSecret   string `env:"ACTION_SK" required:"true"`
	ScalingGroupId string `env:"ACTION_SCALING_GROUP_ID" required:"true"`
	WorkDir        string `env:"WORKDIR" required:"true"`
	client         *api.Client
}

func (e *Ess) Init() error {
	client, err := api.NewClientWithAccessKey(e.Region, e.AccessKeyID, e.AccessSecret)
	if err != nil {
		logrus.Errorf("create ess client error: %v", err)
		return err
	}
	e.client = client
	return nil
}

// check whether scaling group exist
func (e *Ess) IsScaleGroupExist(essGroupName string) (bool, error) {

	request := api.CreateDescribeScalingGroupsRequest()
	request.Scheme = "https"

	request.ScalingGroupName1 = essGroupName

	response, err := e.client.DescribeScalingGroups(request)
	if err != nil {
		logrus.Errorf("failed to get ess group: %v", err)
		return false, err
	}
	if response.TotalCount != 1 {
		return false, nil
	}
	return true, nil
}

// get instance ids in scaling group
func (e *Ess) GetScalingInstances() ([]string, error) {
	var result []string
	var rsp []api.ScalingInstance
	var pageNumber int
	request := api.CreateDescribeScalingInstancesRequest()
	request.Scheme = "https"

	request.ScalingGroupId = e.ScalingGroupId
	pageSize := 10
	request.PageSize = requests.NewInteger(pageSize)

	start := false

	// if it's the first round loop, or return full page size, continue
	for !start || len(rsp) == pageSize {
		start = true
		response, err := e.client.DescribeScalingInstances(request)
		if err != nil {
			// TODO: maybe it's better to continue
			logrus.Errorf("get scaling instances error: %v", err)
			return nil, err
		}
		rsp = response.ScalingInstances.ScalingInstance
		for _, instance := range rsp {
			result = append(result, instance.InstanceId)
		}

		// update to next page number
		pageNumber = response.PageNumber
		request.PageNumber = requests.NewInteger(pageNumber)
	}
	return result, nil
}

func (e *Ess) GetEssInfo() (map[string]string, error) {
	if err := e.Init(); err != nil {
		return nil, err
	}

	instanceIDs, err := e.GetScalingInstances()
	if err != nil {
		return nil, err
	}

	if instanceIDs == nil {
		logrus.Infof("empty scaling group")
		return nil, nil
	}

	result, err := e.GetInstancesPrivateIp(instanceIDs)
	if err != nil {
		return nil, err
	}

	return result, nil
}
