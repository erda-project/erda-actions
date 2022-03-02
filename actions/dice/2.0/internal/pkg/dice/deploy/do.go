package deploy

import (
	"io/ioutil"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda/pkg/retry"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
	"strings"
)

func (d *deploy) Do() (string, map[string]*common.DeployResult, error) {
	// preCheck params
	if err := paramsPreCheck(d.cfg); err != nil {
		return "", nil, err
	}

	// compose request
	req, err := composeRequest(d.cfg)
	if err != nil {
		return "", nil, err
	}

	// crate deployment order with interval
	var resp common.CreateDeploymentOrderResponse
	err = retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).
			Post(d.cfg.DiceOpenapiPrefix).
			Path(common.DeploymentOrderRequestPath).
			Header(common.Authorization, d.cfg.DiceOpenapiToken).JSONBody(&req).Do().JSON(&resp)
		if err != nil {
			logrus.Errorf("failed to create http client, err: %v", err)
			return err
		}

		if !r.IsOK() {
			reqErr := fmt.Errorf("create a deloyment request failed, statusCode: %d, resp:%+v",
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
		return "", nil, err
	}

	// parse deploy result
	ret := make(map[string]*common.DeployResult)

	for appName, info := range resp.Data.Deployments {
		ret[appName] = &common.DeployResult{
			DeploymentId:  info.DeploymentID,
			ApplicationId: info.ApplicationID,
			RuntimeId:     info.RuntimeID,
		}
	}

	return resp.Data.DeploymentOrderId, ret, nil
}

func composeRequest(c *conf.Conf) (*common.CreateDeploymentOrderRequest, error) {
	releaseType := utils.ConvertType(c.ReleaseTye)
	r := &common.CreateDeploymentOrderRequest{
		ReleaseName: c.ReleaseName,
		Workspace:   strings.ToUpper(c.AssignedWorkspace),
		AutoRun:     true,
		Source:      common.SourcePipeline,
		Type:        releaseType,
	}

	switch releaseType {
	case common.TypeApplicationRelease:
		r.ApplicationName = c.ApplicationName
	case common.TypeProjectRelease:
		r.ProjectId = c.ProjectID
	default:
		r.DeployWithoutBranch = c.DeployWithoutBranch
		r.ReleaseId = c.ReleaseID
		if c.ReleaseID == "" {
			releaseID, err := getReleaseId(c.ReleaseIDPath)
			if err != nil {
				logrus.Errorf("failed to get release id: %s", c.ReleaseIDPath)
				r.Print()
				return nil, err
			}
			r.ReleaseId = releaseID
		}
	}

	r.Print()

	return r, nil
}

func getReleaseId(diceHubPath string) (string, error) {
	fileValue, err := ioutil.ReadFile(fmt.Sprint(diceHubPath, "/dicehub_release"))
	if err != nil {
		return "", errors.New("Read file dicehub_release failed.")
	}

	return string(fileValue), nil
}
