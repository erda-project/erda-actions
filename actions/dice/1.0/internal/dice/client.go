package dice

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

type dice struct {
	conf *conf
}

type DiceResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Err     Err         `json:"err,omitempty"`
}

type Err struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"msg,omitempty"`
	Ctx     map[string]interface{} `json:"ctx,omitempty"`
}

type DeployResult struct {
	DeploymentId  int64  `json:"deploymentId"`
	ApplicationId int64  `json:"applicationId"`
	RuntimeId     int64  `json:"runtimeId"`
	Operator      string `json:"operator"`
}

type deployRequest struct {
	ClusterName  string                 `json:"clusterName"`
	Name         string                 `json:"name"`
	Operator     string                 `json:"operator"`
	Source       string                 `json:"source"`
	ReleaseId    string                 `json:"releaseId"`
	EdgeLocation string                 `json:"edgeLocation"`
	Extra        map[string]interface{} `json:"extra, omitempty"`
}

type DiceDeployError struct {
	s string
}

func (e *DiceDeployError) Error() string {
	return e.s
}

const Authorization = "Authorization"

func (d *dice) Deploy(deployReq *deployRequest, conf *conf) (*DeployResult, error) {
	var diceResp DiceResponse
	err := retry.DoWithInterval(func() error {

		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).Path("/api/runtimes").
			Header(Authorization, conf.DiceOpenapiToken).JSONBody(&deployReq).Do().JSON(&diceResp)
		if err != nil {
			return err
		}
		if !r.IsOK() {
			return errors.Errorf("create a dice deploy failed, statusCode: %d, diceResp:%+v",
				r.StatusCode(), diceResp)
		}

		if !diceResp.Success {
			return errors.Errorf("create dice deploy failed. code=%s, message=%s, ctx=%v",
				diceResp.Err.Code, diceResp.Err.Message, diceResp.Err.Ctx)
		}

		return nil
	}, 5, time.Second*3)

	if err != nil {
		logrus.Errorf("deploy to dice failed! response err:%v.", err)
		return nil, err
	}

	result := DeployResult{}
	if err := mapstructure.Decode(diceResp.Data, &result); err != nil {
		return nil, errors.Wrapf(err, "mapstructure data=%+v", result)
	}

	return &result, nil
}

type R struct {
	Success bool `json:"success"`
	Data    struct {
		DeploymentId int               `json:"deploymentId"`
		Status       string            `json:"status"`
		FailCause    string            `json:"failCause"`
		MoudleErrMsg map[string]string `json:"lastMessage"`
		Runtime      interface{}       `json:"runtime"`
	} `json:"data"`
	Err Err `json:"err,omitempty"`
}

func (d *dice) Check(res *DeployResult, conf *conf) (bool, interface{}, error) {
	result, err := getDeploymentStatus(res, conf)
	if err != nil {
		return false, nil, err
	}
	if len(result.Data.MoudleErrMsg) > 0 {
		storeMetaFileWithErr(conf, res.RuntimeId, res.DeploymentId, result)
	}
	switch result.Data.Status {
	case "WAITING", "WAITAPPROVE", "INIT":
		return true, nil, nil
	case "DEPLOYING":
		logrus.Infof("continue deploying..., ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
			res.ApplicationId, res.RuntimeId, res.DeploymentId)
		return true, nil, nil
	case "OK":
		logrus.Info("deploy success")
		return false, result.Data.Runtime, nil
	case "CANCELED":
		return false, nil, &DiceDeployError{"deployment canceled by dice"}
	case "FAILED":
		return false, nil, &DiceDeployError{"deployment failed in dice, " + result.Data.FailCause}
	}
	return false, nil, errors.Errorf("deployment unknown %s in dice", result.Data.Status)
}

func getDeploymentStatus(res *DeployResult, conf *conf) (*R, error) {
	var result R
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path(fmt.Sprintf("/api/deployments/%d/status", res.DeploymentId)).
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&result)
	if err != nil {
		return nil, err
	}
	logrus.Infof("deployments status response body : %v", result)
	if !r.IsOK() {
		return nil, errors.Errorf("deploy to dice failed, statusCode: %d", r.StatusCode())
	}
	if !result.Success {
		return nil, errors.Errorf("create dice deploy failed. code=%s, message=%s, ctx=%v",
			result.Err.Code, result.Err.Message, result.Err.Ctx)
	}
	return &result, nil
}

func getReleaseId(diceHubPath string) (string, error) {
	fileValue, err := ioutil.ReadFile(fmt.Sprint(diceHubPath, "/dicehub_release"))
	if err != nil {
		return "", errors.New("Read file dicehub_release failed.")
	}

	return string(fileValue), nil
}
