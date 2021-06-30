package dice_deploy_services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

type Conf struct {
	OrgID        uint64 `env:"DICE_ORG_ID"`
	ProjectID    uint64 `env:"DICE_PROJECT_ID"`
	AppID        uint64 `env:"DICE_APPLICATION_ID"`
	Workspace    string `env:"DICE_WORKSPACE"`
	GittarBranch string `env:"GITTAR_BRANCH"`
	ClusterName  string `env:"DICE_CLUSTER_NAME"`
	OperatorID   string `env:"DICE_OPERATOR_ID"`

	// used to invoke openapi
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR"`
	DiceOpenapiToken  string `env:"DICE_OPENAPI_TOKEN"`
	InternalClient    string `env:"DICE_INTERNAL_CLIENT"`
	UserID            string `env:"DICE_USER_ID"`

	// wd & meta
	WorkDir  string `env:"WORKDIR"`
	MetaFile string `env:"METAFILE"`

	PipelineBuildID uint64 `env:"PIPELINE_ID"`
	PipelineTaskID  uint64 `env:"PIPELINE_TASK_ID"`

	// params
	DeploymentID string `env:"ACTION_DEPLOYMENT_ID"`
	TimeOut      int    `env:"ACTION_TIME_OUT"`
}

type Dice struct {
	Conf *Conf
}

type Err struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"msg,omitempty"`
	Ctx     map[string]interface{} `json:"ctx,omitempty"`
}

type DiceResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Err     Err         `json:"err,omitempty"`
}
type DeployResult struct {
	DeploymentId  int64 `json:"deploymentId"`
	ApplicationId int64 `json:"applicationId"`
	RuntimeId     int64 `json:"runtimeId"`
}
type DiceDeployError struct {
	S string
}

func (e *DiceDeployError) Error() string {
	return e.S
}

const Authorization = "Authorization"

func Run() error {
	var cfg Conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Info(cfg)
	d := &Dice{Conf: &cfg}

	result, err := Deploy(&cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}

	//Set default deployment timeout is 24h.
	timeout := cfg.TimeOut
	minTimeoutSec := (60 * 60) * 24
	if timeout < minTimeoutSec {
		timeout = minTimeoutSec
	}

	_, err = CheckDeploymentLoop(d, result, "", time.Duration(timeout), Check)
	deployResult, Deployerr := GetDeploymentStatus(result, &cfg)
	if Deployerr != nil {
		return Deployerr
	}
	if err != nil {
		StoreMetaFileWithErr(&cfg, result.RuntimeId, result.DeploymentId, deployResult)
		return err
	}
	logrus.Infof("checkDeploymentLoop end storeMetaFile")
	return StoreMetaFileWithErr(&cfg, result.RuntimeId, result.DeploymentId, deployResult)
}

func Deploy(conf *Conf) (*DeployResult, error) {
	var diceResp DiceResponse
	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).Path(fmt.Sprintf("/api/deployments/%s/actions/deploy-services", conf.DeploymentID)).
			Header(Authorization, conf.DiceOpenapiToken).Do().JSON(&diceResp)
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
		Phase        string            `json:"phase"`
		FailCause    string            `json:"failCause"`
		MoudleErrMsg map[string]string `json:"lastMessage"`
		Runtime      interface{}       `json:"runtime"`
	} `json:"data"`
	Err Err `json:"err,omitempty"`
}

func Check(res *DeployResult, conf *Conf) (deploying bool, runtime interface{}, e error) {
	defer func() {
		if deploying && e == nil {
			Deploy(conf)
		}
	}()
	result, err := GetDeploymentStatus(res, conf)
	if err != nil {
		deploying = false
		runtime = nil
		e = err
		return
	}
	if len(result.Data.MoudleErrMsg) > 0 {
		StoreMetaFileWithErr(conf, res.RuntimeId, res.DeploymentId, result)
	}
	switch result.Data.Status {
	case "WAITING", "WAITAPPROVE", "INIT":
		deploying = true
		runtime = nil
		e = nil
		return
	case "DEPLOYING":
		switch result.Data.Phase {
		case "INIT", "ADDON_REQUESTING", "SCRIPT_APPLYING", "SERVICE_DEPLOYING":
			logrus.Infof("continue deploying..., ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
				res.ApplicationId, res.RuntimeId, res.DeploymentId)
			deploying = true
			runtime = nil
			e = nil
			return
		default:
			logrus.Info("deploy services success")
			deploying = false
			runtime = result.Data.Runtime
			e = nil
			return
		}
		logrus.Infof("continue deploying..., ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
			res.ApplicationId, res.RuntimeId, res.DeploymentId)
		deploying = true
		runtime = nil
		e = nil
		return
	case "OK":
		logrus.Info("deploy success")
		deploying = false
		runtime = result.Data.Runtime
		e = nil
		return
	case "CANCELED":
		deploying = false
		runtime = nil
		e = &DiceDeployError{"deployment canceled by dice"}
		return
	case "FAILED":
		deploying = false
		runtime = nil
		e = &DiceDeployError{"deployment failed in dice, " + result.Data.FailCause}
		return
	}
	deploying = false
	runtime = nil
	e = errors.Errorf("deployment unknown %s in dice", result.Data.Status)
	return
}

// StoreMetaFileWithErr metadata写入err信息
func StoreMetaFileWithErr(conf *Conf, runtimeID int64, deploymentID int64, deployResult *R) error {
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
	logrus.Infof("StoreMetaFileWithErr CreateFile body: %v", string(b))
	if err := filehelper.CreateFile(conf.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func storeMetaFile(conf *Conf, runtimeID int64, deploymentID int64) error {
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
func GetDeploymentStatus(res *DeployResult, conf *Conf) (*R, error) {
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

// generateMetadata 生成固定Metadata数据
func generateMetadata(conf *Conf, runtimeID int64, deploymentID int64) *apistructs.Metadata {
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

func CheckDeploymentLoop(
	d *Dice,
	result *DeployResult,
	operator string,
	timeOut time.Duration,
	checkF func(res *DeployResult, conf *Conf) (bool, interface{}, error),
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
			deploying, runtime, err = checkF(result, d.Conf)
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
