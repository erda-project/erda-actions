package erda_get_addon_info

import (
	"bytes"
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
	AddonName  string `env:"ACTION_ADDON_NAME"`
}

type dice struct {
	conf *conf
}

type Result struct {
	Info   map[string]string   `json:"addonConfigs"`
}

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Infof("Cfg: %#v",cfg)

	d := &dice{conf: &cfg}

	result, err := d.GetAddonInstances(&cfg)
	if err != nil {
		return errors.Wrap(err, "fet addon config info failed")
	}
	storeMetaFile(&cfg, result)
	return nil
}

func (d *dice) GetAddonInstances(conf *conf) (*Result, error) {
	if conf.RuntimeID == "" || conf.AddonName == "" {
		logrus.Errorf("get-addon-info failed: runtime_id and addon_name must provided.")
		return nil, errors.Errorf("get-addon-info failed: runtime_id and addon_name must provided.")
	}
	result := &Result{}

	err := retry.DoWithInterval(func() error {
		result.Info = make(map[string]string)
		addons, err := getAddons(conf)
		if err != nil {
			logrus.Errorf("get-addon-info call getAddons failed, error: %v", err)
			return err
		}

		addonId := ""
		for _, addon := range addons {
			if addon.Name == conf.AddonName {
				addonId = addon.ID
				break
			}
		}

		if addonId == "" {
			logrus.Infof("get-addon-info successfully, but no addon found for runtimeID %s with name %s", conf.RuntimeID, conf.AddonName)
			return nil
		}

		if addonId != "" {
			addonInfo, err := getAddonByRoutingKeyId(conf, addonId)
			if err != nil {
				logrus.Errorf("get-addon-info call getAddons failed, error: %v", err)
				return err
			}

			if addonInfo.ID == "" {
				logrus.Errorf("get-addon-info call getAddons failed, no addon found for addon routing key id %s.", addonId)
				return errors.Errorf("get-addon-info call getAddons failed, no addon found for addon routing key id %s.", addonId)
			}

			if addonInfo.Config == nil {
				logrus.Infof("get-addon-info successfully, but no addon config found for runtimeID %s with name %s", conf.RuntimeID, conf.AddonName)
				return nil
			}

			for key, value := range addonInfo.Config {
				result.Info[key] = fmt.Sprintf("%v", value)
			}
		}

		return nil
	}, 5, time.Second*3)
	if err != nil {
		logrus.Errorf("get-addon-info failed! response err:%v.", err)
		return nil, err
	}

	return result, nil
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
	for name, addr := range result.Info {
		addrs = append(addrs, apistructs.MetadataField{
			Name:  name,
			Value: addr,
		} )
	}

	var ret apistructs.Metadata
	ret = addrs
	return &ret
}


// 获取应用程序对应的 Runtime
func getAddons(conf *conf) ([]apistructs.AddonFetchResponseData, error){
	var resp apistructs.AddonListResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path("/api/addons").
		Param("type", "runtime").
		Param("workspace", conf.Workspace).
		Param("value", conf.RuntimeID).
		Param("projectId", fmt.Sprintf("%d",conf.ProjectID)).
		Header("User-ID", conf.UserID).
		Header("Org-ID", strconv.FormatUint(conf.OrgID, 10)).
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)

	if err != nil {
		logrus.Infof("getAddons for runtimeID %s failed, error: %v", conf.RuntimeID, err)
		return nil, err
	}

	if !r.IsOK() || !resp.Success {
		logrus.Infof("getAddons for runtimeID %s failed, error msg: %v", conf.RuntimeID, fmt.Errorf(resp.Error.Msg))
		return nil, fmt.Errorf(resp.Error.Msg)
	}
	if len(resp.Data) == 0  {
		logrus.Infof("not found addons for runtimeID %s, error: %v", conf.RuntimeID, fmt.Errorf("runtime ID not found"))
		return nil, fmt.Errorf("addons for runtimeID %s not found", conf.RuntimeID)
	}
	return resp.Data, nil
}


func getAddonByRoutingKeyId(cfg *conf, addonID string) (*apistructs.AddonFetchResponseData, error) {
	var buffer bytes.Buffer
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/addons/%s", addonID)).
		Header("Authorization", cfg.DiceOpenapiToken).
		Do().Do().Body(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to getAddonFetchResponseData, err: %v", err)
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("failed to getAddonFetchResponseData, statusCode: %d, respBody: %s", resp.StatusCode(), buffer.String())
	}
	var result apistructs.AddonFetchResponse
	respBody := buffer.String()
	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return nil, fmt.Errorf("failed to getAddonFetchResponseData, err: %v, json string: %s", err, respBody)
	}
	return &result.Data, nil
}