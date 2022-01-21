package erda_get_service_addr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	RuntimeID string `env:"ACTION_RUNTIME_ID"`
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

type Result struct {
	ServicesAddrs   map[string]string `json:"services"`
}

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Infof("%#v",cfg)

	d := &dice{conf: &cfg}

	result, err := d.GetServices(&cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}
	storeMetaFile(&cfg, result)
	return nil
}

func (d *dice) GetServices(conf *conf) (*Result, error) {
	var resp apistructs.RuntimeInspectResponse
	if conf.RuntimeID == "" {
		logrus.Errorf("GetServices failed: runtimeID not provided.")
		return nil, errors.Errorf("GetServices failed: runtimeID not provided.")
	}

	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path(fmt.Sprintf("/api/runtimes/%s", conf.RuntimeID)).
			Param("applicationId", fmt.Sprintf("%d",conf.ProjectID)).
			Param("workspace", conf.Workspace).
			Header("User-ID", conf.UserID).
			Header("Org-ID", strconv.FormatUint(conf.OrgID, 10)).
			Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)

		if err != nil {
			logrus.Errorf("get Runtime failed, error: %v", err)
			return err
		}
		if !r.IsOK() {
			logrus.Errorf("get Runtime failed, statusCode: %v, resp:%+v", resp.Error.Code, resp)
			return errors.Errorf("get Runtime failed, statusCode: %v, resp:%+v", resp.Error.Code, resp)
		}

		if !resp.Success {
			logrus.Errorf("get Runtime failed. code=%s, message=%s, ctx=%v", resp.Error.Code, resp.Error.Msg, resp.Error.Ctx)
			return errors.Errorf("get Runtime failed. code=%s, message=%s, ctx=%v", resp.Error.Code, resp.Error.Msg, resp.Error.Ctx)
		}
		return nil
	}, 5, time.Second*3)
	if err != nil {
		logrus.Errorf("get runtime services failed! response err:%v.", err)
		return nil, err
	}

	if resp.Data.ID == 0 || len(resp.Data.Services) == 0  {
		logrus.Infof("not found runtime for runtime Id %s , error: %v", conf.RuntimeID, fmt.Errorf("runtime not found"))
		return nil, fmt.Errorf("runtime not found")
	}
	result := Result{
		ServicesAddrs: make(map[string]string),
	}

	for name, svc := range resp.Data.Services {
		if len(svc.Addrs) > 0 {
			result.ServicesAddrs[name]= svc.Addrs[0]
			continue
		}
	}

	logrus.Infof("get Runtime %s services sucessfully. addrs: %#v", conf.RuntimeID, result.ServicesAddrs)
	return &result, nil
}

func storeMetaFile(conf *conf, result *Result) error {
	meta := apistructs.ActionCallback{
		Metadata: *generateMetadata(result),
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

func generateMetadata(result *Result) *apistructs.Metadata {
	addrs := make([]apistructs.MetadataField, 0)
	for name, addr := range result.ServicesAddrs {
		addrs = append(addrs, apistructs.MetadataField{
			Name:  name,
			Value: addr,
		} )
	}

	var ret apistructs.Metadata
	ret = addrs
	return &ret
}
