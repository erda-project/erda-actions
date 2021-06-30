package dice_deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ReleaseID     string `env:"ACTION_RELEASE_ID"`
	ReleaseIDPath string `env:"ACTION_RELEASE_ID_PATH"`
}
type dice struct {
	conf *conf
}

type deployRequest struct {
	ClusterName    string                 `json:"clusterName"`
	Name           string                 `json:"name"`
	Operator       string                 `json:"operator"`
	Source         string                 `json:"source"`
	ReleaseId      string                 `json:"releaseId"`
	Extra          map[string]interface{} `json:"extra, omitempty"`
	SkipPushByOrch bool                   `json:"skipPushByOrch"`
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

func prepareRequest(conf *conf) (*deployRequest, error) {
	req := new(deployRequest)
	req.ClusterName = conf.ClusterName
	req.Name = conf.GittarBranch
	req.Operator = conf.OperatorID
	req.Source = "PIPELINE"
	req.SkipPushByOrch = true

	extra := make(map[string]interface{})
	extra["orgId"] = int(conf.OrgID)
	extra["projectId"] = int(conf.ProjectID)
	extra["applicationId"] = int(conf.AppID)
	extra["workspace"] = conf.Workspace
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

func getReleaseId(diceHubPath string) (string, error) {
	fileValue, err := ioutil.ReadFile(fmt.Sprint(diceHubPath, "/dicehub_release"))
	if err != nil {
		return "", errors.New("Read file dicehub_release failed.")
	}

	return string(fileValue), nil
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

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Info(cfg)
	d := &dice{conf: &cfg}

	deployReq, err := prepareRequest(&cfg)
	if err != nil {
		return errors.Wrap(err, "prepare dice deploy request failed")
	}
	result, err := d.Deploy(deployReq, &cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}
	storeMetaFile(&cfg, result.RuntimeId, result.DeploymentId)
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
