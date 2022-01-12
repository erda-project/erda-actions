package dice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

type dice struct {
	conf *conf
}

func (d *dice) Deploy(deployReq *CreateDeploymentOrderRequest, conf *conf) (*DeployResult, error) {
	var resp CreateDeploymentOrderResponse
	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).Path(DeploymentOrderRequestPath).
			Header(Authorization, conf.DiceOpenapiToken).JSONBody(&deployReq).Do().JSON(&resp)
		if err != nil {
			logrus.Errorf("failed to create http client, err: %v", err)
			return err
		}

		if !r.IsOK() {
			reqErr := fmt.Errorf("create a dice deploy failed, statusCode: %d, resp:%+v",
				r.StatusCode(), resp)
			logrus.Error(reqErr)
			return reqErr
		}

		if !resp.Success {
			respErr := errors.Errorf("create dice deploy failed. code=%s, message=%s, ctx=%v",
				resp.Err.Code, resp.Err.Message, resp.Err.Ctx)
			logrus.Error(respErr)
			return respErr
		}

		logrus.Infof("request response: %+v", resp)

		return nil
	}, 5, time.Second*3)

	if err != nil {
		logrus.Errorf("deploy to dice failed! response err: %v", err)
		return nil, err
	}

	deployment, ok := resp.Data.Deployments[conf.AppID]
	if !ok {
		return nil, errors.Wrapf(err, "get deployment info from reponse error, application: %d", conf.AppID)
	}

	return &DeployResult{
		DeploymentId:  deployment.DeploymentID,
		ApplicationId: deployment.ApplicationID,
		RuntimeId:     deployment.RuntimeID,
	}, nil
}

func (d *dice) Check(res *DeployResult, conf *conf, lastDeployStatusInfo *string) (bool, interface{}, error) {
	result, err := getDeploymentStatus(res, conf)
	if err != nil {
		return false, nil, err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return false, nil, err
	}
	deployStatusInfo := string(b)
	if deployStatusInfo != *lastDeployStatusInfo {
		*lastDeployStatusInfo = deployStatusInfo
		result.Print()
	}

	if len(result.Data.ModuleErrMsg) > 0 {
		storeMetaFileWithErr(conf, res.RuntimeId, res.DeploymentId, result)
	}
	switch result.Data.Status {
	case "WAITING", "WAITAPPROVE", "INIT":
		return true, nil, nil
	case "DEPLOYING":
		return true, nil, nil
	case "OK":
		logrus.Info("deploy success!")
		return false, result.Data.Runtime, nil
	case "CANCELED":
		return false, nil, &DeployErrResponse{"deployment canceled by dice"}
	case "FAILED":
		return false, nil, &DeployErrResponse{"deployment failed in dice, " + result.Data.FailCause}
	}
	return false, nil, errors.Errorf("deployment unknown %s in dice", result.Data.Status)
}

func getDeploymentStatus(res *DeployResult, conf *conf) (*DeploymentStatusRespData, error) {
	var result DeploymentStatusRespData
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path(fmt.Sprintf("/api/deployments/%d/status", res.DeploymentId)).
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&result)
	if err != nil {
		return nil, err
	}
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
