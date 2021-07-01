package dice_deploy_release

import (
	"encoding/json"
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
	ReleaseID         string `env:"ACTION_RELEASE_ID"`
	AssignedWorkspace string `env:"ACTION_WORKSPACE"`
}
type dice struct {
	conf *conf
}

type deployRequest struct {
	// 制品ID
	ReleaseID string `json:"releaseId"`
	// 环境
	Workspace string `json:"workspace"`
	// 项目ID
	ProjectID uint64 `json:"projectId"`
	// 应用ID
	ApplicationID uint64 `json:"applicationId"`
}

func prepareRequest(conf *conf) *deployRequest {
	req := new(deployRequest)
	req.ReleaseID = conf.ReleaseID
	req.Workspace = conf.Workspace
	if conf.AssignedWorkspace != "" {
		req.Workspace = conf.AssignedWorkspace
	}
	req.ProjectID = conf.ProjectID
	req.ApplicationID = conf.AppID
	return req
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

func (d *dice) Deploy(deployReq *deployRequest, conf *conf) (*DeployResult, error) {
	var diceResp DiceResponse
	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).
			Path("/api/runtimes/actions/deploy-release-action").
			Header("Authorization", conf.DiceOpenapiToken).JSONBody(&deployReq).Do().JSON(&diceResp)
		if err != nil {
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

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Info(cfg)
	d := &dice{conf: &cfg}

	deployReq := prepareRequest(&cfg)
	result, err := d.Deploy(deployReq, &cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}
	storeMetaFile(&cfg, int64(result.RuntimeID), int64(result.DeploymentID))
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
