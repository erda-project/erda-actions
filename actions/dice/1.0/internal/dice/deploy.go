package dice

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/http/httputil"
	"github.com/erda-project/erda/pkg/retry"
)

func Run() error {
	cfg, err := HandleConf()
	if err != nil {
		return err
	}
	d := &dice{conf: &cfg}

	deployReq, err := prepareRequest(&cfg, cfg.AssignedWorkspace)
	if err != nil {
		return errors.Wrap(err, "prepare dice deploy request failed")
	}

	result, err := d.Deploy(deployReq, &cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}
	//store dice deploymentId and runtimeID.
	runtimeID := strconv.FormatInt(result.RuntimeId, 10)
	deploymentId := strconv.FormatInt(result.DeploymentId, 10)
	err = storeDiceInfo(deploymentId, runtimeID, d.conf.WorkDir)
	if err != nil {
		logrus.Warning(err)
	}
	// Report runtimeID to pipeline platform
	if runtimeID == "" {
		logrus.Warningf("runtimeID is 0. can not report to ci")
	} else {
		err := reportRuntimeID2PipelinePlatform(&cfg, runtimeID)
		if err != nil {
			logrus.Warningf("Report runtimeID to ci failed.")
		}
	}

	//Set default deployment timeout is 24h.
	timeout := cfg.TimeOut
	minTimeoutSec := (60 * 60) * 24
	if timeout < minTimeoutSec {
		timeout = minTimeoutSec
	}
	runtime, err := checkDeploymentLoop(
		d,
		result,
		fmt.Sprintf("%v", cfg.OperatorID),
		time.Duration(timeout),
	)
	deployResult, Deployerr := getDeploymentStatus(result, &cfg)
	if Deployerr != nil {
		return Deployerr
	}

	if err != nil {
		storeMetaFileWithErr(&cfg, result.RuntimeId, result.DeploymentId, deployResult)
		callback(cfg.Callback, 0, cfg.AppID, runtime, "Failed")
		return err
	} else {
		// do callback
		callback(cfg.Callback, result.RuntimeId, cfg.AppID, runtime, "Success")
	}
	logrus.Infof("checkDeploymentLoop end storeMetaFile")
	return storeMetaFileWithErr(&cfg, result.RuntimeId, result.DeploymentId, deployResult)
}

func storeDiceInfo(deploymentId, runtimeId, wd string) error {
	content := fmt.Sprint("deploymentId=", deploymentId, ",", "runtimeId=", runtimeId)
	err := filehelper.CreateFile(filepath.Join(wd, "diceInfo"), content, 0755)
	if err != nil {
		return errors.Wrap(err, "write file:diceInfo failed")
	}
	return nil
}

// generateMetadata 生成固定Metadata数据
func generateMetadata(conf *conf, runtimeID int64, deploymentID int64) *apistructs.Metadata {
	return &apistructs.Metadata{
		{
			Name:  "project_id",
			Value: strconv.FormatUint(conf.ProjectID, 10),
		},
		{
			Name:  "app_id",
			Value: strconv.FormatUint(conf.AppID, 10),
		},
		{
			Name:  apistructs.ActionCallbackRuntimeID,
			Value: strconv.FormatInt(runtimeID, 10),
			Type:  apistructs.ActionCallbackTypeLink,
		},
		{
			Name:  "deployment_id",
			Value: strconv.FormatInt(deploymentID, 10),
		},
	}
}

// storeMetaFileWithErr metadata写入err信息
func storeMetaFileWithErr(conf *conf, runtimeID int64, deploymentID int64, deployResult *R) error {
	if deployResult == nil {
		return storeMetaFile(conf, runtimeID, deploymentID)
	}
	if len(deployResult.Data.MoudleErrMsg) == 0 {
		return storeMetaFile(conf, runtimeID, deploymentID)
	}
	metadata := generateMetadata(conf, runtimeID, deploymentID)
	for k, v := range deployResult.Data.MoudleErrMsg {
		*metadata = append(*metadata, apistructs.MetadataField{
			Name:  k,
			Value: v,
		})
	}
	meta := apistructs.ActionCallback{
		Metadata: *metadata,
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	logrus.Infof("storeMetaFileWithErr CreateFile body: %v", string(b))
	if err := filehelper.CreateFile(conf.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func storeMetaFile(conf *conf, runtimeID int64, deploymentID int64) error {
	meta := apistructs.ActionCallback{
		Metadata: *generateMetadata(conf, runtimeID, deploymentID),
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(conf.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func reportRuntimeID2PipelinePlatform(conf *conf, runtimeID string) error {
	cb := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{Name: apistructs.ActionCallbackRuntimeID, Value: runtimeID, Type: apistructs.ActionCallbackTypeLink},
			{Name: apistructs.ActionCallbackOperatorID, Value: conf.OperatorID},
		},
		PipelineID:     conf.PipelineBuildID,
		PipelineTaskID: conf.PipelineTaskID,
	}
	b, err := json.Marshal(&cb)
	if err != nil {
		return err
	}

	var cbReq apistructs.PipelineCallbackRequest
	cbReq.Type = string(apistructs.PipelineCallbackTypeOfAction)
	cbReq.Data = b

	if err = retry.DoWithInterval(func() error {
		var resp apistructs.PipelineCallbackResponse
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).
			Post(conf.DiceOpenapiPrefix).
			Path("/api/pipelines/actions/callback").
			Header("Authorization", conf.DiceOpenapiToken).
			Header(httputil.UserHeader, conf.UserID).
			Header(httputil.InternalHeader, conf.InternalClient).
			JSONBody(&cbReq).
			Do().
			JSON(&resp)
		if err != nil {
			return err
		}
		if !r.IsOK() || !resp.Success {
			return errors.Errorf("status-code %d, resp %#v", r.StatusCode(), resp)
		}
		logrus.Infof("report runtimeID to pipeline platform successfully! runtimeID: %s", runtimeID)
		return nil
	}, 3, time.Millisecond*500); err != nil {
		logrus.Infof("report runtimeID to pipeline platform failed! err: %v", err)
		return err
	}
	return nil
}

func checkDeploymentLoop(
	d *dice,
	result *DeployResult,
	operator string,
	timeOut time.Duration,
) (interface{}, error) {
	timer := time.NewTimer(timeOut * time.Second)
	deploying := true
	var runtime interface{}

	// Check if APP was deployed.
deployloop:
	for {
		select {
		case <-timer.C:
			break deployloop
		default:
			var err error
			deploying, runtime, err = d.Check(result, d.conf)
			if err != nil {
				logrus.Errorf("check deploying is not null")
				if _, ok := err.(*DiceDeployError); ok {
					logrus.Errorf("Deploy to Dice Failed: %s", err.Error())
					logrus.Errorf("Deployment link ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
						result.ApplicationId, result.RuntimeId, result.DeploymentId)
				}
				return nil, err
			}
			if !deploying {
				break deployloop
			}
		}

		time.Sleep(10 * time.Second)
	}
	logrus.Errorf("deployloop continue")
	if deploying {
		logrus.Errorf("Deploying timeout( %d seconds). you can: ", timeOut)
		logrus.Error("   1. increase timeout in pipeline.yml")
		logrus.Error("   2. try again ")
		logrus.Errorf("Getting deployment logs ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
			result.ApplicationId, result.RuntimeId, result.DeploymentId)
		//logrus.Error("Now we are going to cancel the task...")
		//cReq := &cancelReq{
		//	DeploymentId: result.DeploymentId,
		//	RuntimeId:    result.RuntimeId,
		//	Operator:     operator,
		//}
		//err := d.Cancel(cReq, envs)
		//if err != nil {
		//	return nil, errors.Wrapf(err, "cancel deployment with req=%v failed", cReq)
		//}
		//return nil, errors.New("deployment canceled")
		return nil, errors.New("deployment timeout")
	}
	logrus.Errorf("return runtime")
	return runtime, nil
}

func prepareRequest(conf *conf, workspace string) (*deployRequest, error) {
	req := new(deployRequest)
	req.ClusterName = conf.ClusterName
	req.Name = conf.GittarBranch
	req.Operator = conf.OperatorID
	req.Source = "PIPELINE"
	req.EdgeLocation = conf.EdgeLocation

	extra := make(map[string]interface{})
	extra["orgId"] = int(conf.OrgID)
	extra["projectId"] = int(conf.ProjectID)
	extra["applicationId"] = int(conf.AppID)
	extra["workspace"] = conf.Workspace
	if workspace != "" {
		extra["workspace"] = workspace
	}
	extra["buildId"] = conf.PipelineBuildID

	logrus.Infof("<<<request deploy body:%v", req)

	req.Extra = extra

	var releaseID string
	if conf.ReleaseID != "" {
		releaseID = conf.ReleaseID
	} else {
		var err error
		releaseID, err = getReleaseId(conf.ReleaseIDPath)
		if err != nil {
			return nil, err
		}
	}

	logrus.Infof("<<<releaseID:%s", releaseID)

	req.ReleaseId = releaseID

	return req, nil
}

func callback(url string, runtimeId int64, applicationId uint64, options interface{}, status string) {
	if len(url) == 0 {
		logrus.Info("no callback set, return directly")
		return
	}
	logrus.Infof("start to callback with url=%s, runtimeId=%d, applicationId=%d, status=%s, options=%v ...", url, runtimeId, applicationId, status, options)
	data := struct {
		ApplicationId uint64      `json:"applicationId"`
		RuntimeId     int64       `json:"runtimeId"`
		Status        string      `json:"status"`
		Options       interface{} `json:"options, omitempty"`
	}{
		ApplicationId: applicationId,
		RuntimeId:     runtimeId,
		Status:        status,
		Options:       options,
	}
	dj, err := json.Marshal(&data)
	if err != nil {
		logrus.Error(err)
	}
	resp, body, errs := gorequest.New().Post(url).Send(string(dj)).Timeout(5 * time.Second).End()

	if errs != nil {
		logrus.Error(multierror.Append(errors.New("callback failed"), errs...))
	}
	if resp.StatusCode/100 != 2 {
		logrus.Error(errors.Errorf("callback failed, response code: %d, response body: %s",
			resp.StatusCode, body))
	}
	logrus.Infof("callback successfully! response code=%d", resp.StatusCode)
}
