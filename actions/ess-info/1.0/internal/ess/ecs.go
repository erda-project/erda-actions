package ess

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func (e *Ess) GetInstancesInfo(instanceIds []string) (*ecs.DescribeInstancesResponse, error) {
	client, err := ecs.NewClientWithAccessKey(e.Region, e.AccessKeyID, e.AccessSecret)
	if err != nil {
		return nil, err
	}
	if instanceIds == nil {
		return nil, errors.New("empty instance ids")
	}

	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.RegionId = e.Region

	content, err := json.Marshal(instanceIds)
	if err != nil {
		errStr := fmt.Sprintf("json marshal error: %v", err)
		return nil, errors.New(errStr)
	}
	request.InstanceIds = string(content)
	request.PageSize = requests.Integer(strconv.Itoa(len(instanceIds)))

	// if not provide valid instance id, it will return other instance info by default pageNum & pageSize
	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}
	if response.BaseResponse == nil {
		return nil, errors.New("base response in empty")
	}

	if !response.BaseResponse.IsSuccess() {
		errStr := fmt.Sprintf("base response status code: %d", response.BaseResponse.GetHttpStatus())
		return nil, errors.New(errStr)
	}

	return response, nil
}

func (e *Ess) GetInstancesPrivateIp(instanceIds []string) (map[string]string, error) {
	result := make(map[string]string)
	response, err := e.GetInstancesInfo(instanceIds)
	if err != nil {
		return nil, err
	}
	validNum := 0
	for _, instance := range response.Instances.Instance {
		if instance.VpcAttributes.PrivateIpAddress.IpAddress == nil {
			return nil, errors.New("get empty instance private ip")
		}
		exist, err := contains(instanceIds, instance.InstanceId)
		if err != nil {
			return nil, err
		}
		if !exist {
			errStr := fmt.Sprintf("instance id: %s, not in request instance ids: %v", instance.InstanceId, instanceIds)
			return nil, errors.New(errStr)
		}
		result[instance.InstanceId] = instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
		validNum++
	}
	if validNum != len(instanceIds) {
		errStr := fmt.Sprintf("valid instance num: %d, total num: %d, response: %v, all instance ids: %v", validNum, len(instanceIds), result, instanceIds)
		return nil, errors.New(errStr)
	}
	return result, nil

}

func contains(slice interface{}, item interface{}) (bool, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return false, errors.New("invalid data type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true, nil
		}
	}
	return false, nil
}
