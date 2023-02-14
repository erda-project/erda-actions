package erda_add_app_config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/pkg/metadata"
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
	ApplicationName string            `env:"ACTION_APPLICATION_NAME"`
	ConfigWorkspace string            `env:"ACTION_CONFIG_WORKSPACE"`
	ConfigItems     map[string]string `env:"ACTION_CONFIG_ITEMS"`
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

type GetAppIDByNameResult map[string]int64

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Info(cfg)
	d := &dice{conf: &cfg}

	result, err := d.GetAppIDByName(&cfg)
	if err != nil {
		return errors.Wrap(err, "get appID by appName failed")
	}
	err = d.AddConfig(&cfg, result)
	if err != nil {
		return errors.Wrap(err, "add config for application failed")
	}
	storeMetaFile(&cfg)
	return nil
}

func (d *dice) GetAppIDByName(conf *conf) (uint64, error) {
	var resp apistructs.ApplicationListResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path("/api/applications").
		Param("projectId", fmt.Sprintf("%d", conf.ProjectID)).
		Param("name", conf.ApplicationName).
		Param("pageNo", "1").
		Param("pageSize", "1").
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)

	if err != nil {
		logrus.Infof("getAppID for app name %s failed, error: %v", conf.ApplicationName, err)
		return 0, err
	}
	if !r.IsOK() || !resp.Success {
		return 0, errors.Errorf("getAppID for app name %s failed, statusCode: %d, diceResp:%+v", conf.ApplicationName,
			r.StatusCode(), resp)
	}
	if resp.Data.Total == 0 || len(resp.Data.List) == 0 {
		logrus.Infof("not found app for name %s error: %v", conf.ApplicationName, fmt.Errorf("application not found"))
		return 0, fmt.Errorf("application not found")
	}
	return resp.Data.List[0].ID, nil
}

func (d *dice) AddConfig(conf *conf, appID uint64) error {
	configNamespace := fmt.Sprintf("app-%d-%s", appID, strings.ToUpper(conf.ConfigWorkspace))
	var req apistructs.EnvConfigAddOrUpdateRequest
	for k, v := range conf.ConfigItems {
		req.Configs = append(req.Configs, apistructs.EnvConfig{
			Key:        k,
			Value:      v,
			ConfigType: "kv",
			Encrypt:    false,
		})
	}
	var diceResp DiceResponse
	url := fmt.Sprintf("/api/configmanage/configs?namespace_name=%s&encrypt=false&appID=%d", configNamespace, appID)
	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).
			Path(url).
			Header("Authorization", conf.DiceOpenapiToken).JSONBody(req).Do().JSON(&diceResp)
		if err != nil {
			return err
		}
		if !r.IsOK() {
			return errors.Errorf("add config for application failed, statusCode: %d, diceResp:%+v",
				r.StatusCode(), diceResp)
		}

		if !diceResp.Success {
			return errors.Errorf("add config for application failed. code=%s, message=%s, ctx=%v",
				diceResp.Err.Code, diceResp.Err.Message, diceResp.Err.Ctx)
		}
		return nil
	}, 5, time.Second*3)
	if err != nil {
		logrus.Errorf("add config for application! response err:%v.", err)
		return err
	}
	return nil
}

func storeMetaFile(conf *conf) error {
	meta := apistructs.ActionCallback{
		Metadata: *generateMetadata(conf),
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

func generateMetadata(conf *conf) *metadata.Metadata {
	return &metadata.Metadata{
		{
			Name:  "project_id",
			Value: strconv.FormatUint(conf.ProjectID, 10),
		},
	}
}
