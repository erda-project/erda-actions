package dice_deploy_redeploy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/erda-project/erda/pkg/metadata"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

type conf struct {
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
	RuntimeID       string `env:"ACTION_RUNTIME_ID"`
	ApplicationName string `env:"ACTION_APPLICATION_NAME"`
}

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
	DeploymentID  uint64 `json:"deploymentId"`
	ApplicationID uint64 `json:"applicationId"`
	RuntimeID     uint64 `json:"runtimeId"`
}

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Infof("%#v", cfg)

	d := &dice{conf: &cfg}

	result, err := d.Deploy(&cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}
	storeMetaFile(&cfg, int64(result.RuntimeID), int64(result.DeploymentID))
	return nil
}

func (d *dice) Deploy(conf *conf) (*DeployResult, error) {
	var diceResp DiceResponse
	if conf.RuntimeID == "" && conf.ApplicationName == "" {
		logrus.Errorf("deploy failed: neither runtimeID nor ApplicationName provided.")
		return nil, errors.Errorf("deploy failed: neither runtimeID nor ApplicationName provided.")
	}
	if conf.ApplicationName != "" {
		// 通过 ApplicationName 的方式仅支持 trantor 业务以及通 release 部署，走 pipeline 部署的还是需要提供 RuntimeID
		// 通过 application name 先获取到 Application ID，然后结合 WorkSpace 和  应用程序名称（通过 release 部署的） 获取到 runtimeID
		appId, err := getAppID(conf, conf.ApplicationName)
		if err != nil {
			logrus.Errorf("deploy failed: get app Id for appName %s failed, error: %v", conf.ApplicationName, err)
			return nil, errors.Errorf("deploy failed: get app Id for appName %s failed, error: %v.", conf.ApplicationName, err)
		}

		runtimeId, err := getRuntimeId(conf, conf.ApplicationName, appId)
		if err != nil {
			logrus.Errorf("deploy failed: get runtime ID for appName %s failed, error: %v", conf.ApplicationName, err)
			return nil, errors.Errorf("deploy failed: get runtime ID for appName %s failed, error: %v", conf.ApplicationName, err)
		}

		conf.RuntimeID = runtimeId
	}

	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).
			Path(fmt.Sprintf("/api/runtimes/%s/actions/redeploy-action", conf.RuntimeID)).
			Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&diceResp)
		if err != nil {
			logrus.Errorf("redeploy-action failed, error: %v", err)
			return err
		}
		if !r.IsOK() {
			return errors.Errorf("create a dice release deploy failed, statusCode: %d, diceResp:%+v",
				r.StatusCode(), diceResp)
		}

		if !diceResp.Success {
			return errors.Errorf("create dice release deploy failed. code=%s, message=%s, ctx=%v",
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

func generateMetadata(conf *conf, runtimeID int64, deploymentID int64) *metadata.Metadata {
	return &metadata.Metadata{
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

// 获取应用程序 ID
func getAppID(conf *conf, name string) (uint64, error) {
	var resp apistructs.ApplicationListResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path("/api/applications").
		Param("projectId", fmt.Sprintf("%d", conf.ProjectID)).
		Param("name", name).
		Param("pageNo", "1").
		Param("pageSize", "1").
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)

	if err != nil {
		logrus.Infof("getAppID for app name %s failed, error: %v", name, err)
		return 0, err
	}
	if !r.IsOK() || !resp.Success {
		logrus.Infof("getAppID for app name %s failed, error msg: %v", name, fmt.Errorf(resp.Error.Msg))
		return 0, fmt.Errorf(resp.Error.Msg)
	}
	if resp.Data.Total == 0 || len(resp.Data.List) == 0 {
		logrus.Infof("not found app for name %s error: %v", name, fmt.Errorf("application not found"))
		return 0, fmt.Errorf("application not found")
	}
	return resp.Data.List[0].ID, nil
}

// 获取应用程序对应的 Runtime
func getRuntimeId(conf *conf, name string, appId uint64) (string, error) {
	var resp apistructs.RuntimeListResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path("/api/runtimes").
		Param("projectId", fmt.Sprintf("%d", conf.ProjectID)).
		Param("applicationId", fmt.Sprintf("%d", appId)).
		Param("workspace", conf.Workspace).
		Param("name", name).
		Header("User-ID", conf.UserID).
		Header("Org-ID", strconv.FormatUint(conf.OrgID, 10)).
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)

	if err != nil {
		logrus.Infof("getRuntimeId for app Name: %s Id: %d, error: %v", name, appId, err)
		return "", err
	}

	if !r.IsOK() || !resp.Success {
		logrus.Infof("getRuntimeId for app Name: %s Id: %d, error msg: %v", name, appId, fmt.Errorf(resp.Error.Msg))
		return "", fmt.Errorf(resp.Error.Msg)
	}
	if len(resp.Data) == 0 || resp.Data[0].ID == 0 {
		logrus.Infof("not found runtime id for app Name: %s Id: %d, error: %v", name, appId, fmt.Errorf("runtime ID not found"))
		return "", fmt.Errorf("runtime ID not found")
	}
	return strconv.FormatUint(resp.Data[0].ID, 10), nil
}
