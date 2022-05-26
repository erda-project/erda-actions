package deploy

import (
	"io/ioutil"
	"fmt"
	"time"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda/pkg/retry"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func (d *deploy) Do() (string, map[string]*common.DeployResult, error) {
	// preCheck params
	if err := paramsPreCheck(d.cfg); err != nil {
		_ = metawriter.WriteError(err.Error())
		return "", nil, err
	}

	// compose request
	req, err := composeRequest(d.cfg)
	if err != nil {
		return "", nil, err
	}

	var deployError error

	// crate deployment order with interval
	var resp common.CreateDeploymentOrderResponse
	if err = retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).
			Post(d.cfg.DiceOpenapiPrefix).
			Path(common.DeploymentOrderRequestPath).
			Header(common.Authorization, d.cfg.DiceOpenapiToken).JSONBody(&req).Do().JSON(&resp)
		if err != nil {
			return fmt.Errorf("failed to create http client, err: %v", err)
		}

		if !resp.Success {
			respErrs := []string{
				"", // empty line
				fmt.Sprintf("status code: %d", r.StatusCode()),
				fmt.Sprintf("response code: %s", resp.Err.Code),
				fmt.Sprintf("message: %s", resp.Err.Message),
				fmt.Sprintf("context: %s", resp.Err.Ctx),
			}
			respErr := errors.New(strings.Join(respErrs, "\n"))
			// retry
			if resp.Err.Ctx["deploymentOrderID"] == "" {
				return respErr
			}
			// if deployment order already created, break
			deployError = respErr
			_ = metawriter.WriteError(resp.Err.Message)
			return nil
		}

		return nil
	}, 3, time.Second*5); err != nil {
		return "", nil, err
	}

	if deployError != nil {
		return "", nil, deployError
	}

	logrus.Infof("request response: %+v", resp)

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
		Workspace:   c.Workspace,
		AutoRun:     true,
		Source:      common.SourcePipeline,
		Type:        releaseType,
	}

	if c.AssignedWorkspace != "" {
		r.Workspace = c.AssignedWorkspace
	}

	r.Workspace = strings.ToUpper(r.Workspace)

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
