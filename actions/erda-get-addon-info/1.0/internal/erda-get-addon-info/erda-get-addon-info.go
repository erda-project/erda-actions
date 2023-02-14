package erda_get_addon_info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
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
	RuntimeID       string `env:"ACTION_RUNTIME_ID"`
	AddonName       string `env:"ACTION_ADDON_NAME"`
	ApplicationName string `env:"ACTION_APPLICATION_NAME"`
}

type dice struct {
	conf *conf
}

type Result struct {
	Info map[string]string `json:"addonConfigs"`
}

func Run() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Infof("Cfg: %#v", cfg)

	d := &dice{conf: &cfg}

	result, err := d.GetAddonInstances(&cfg)
	if err != nil {
		return errors.Wrap(err, "fet addon config info failed")
	}
	storeMetaFile(&cfg, result)
	return nil
}

func (d *dice) GetAddonInstances(conf *conf) (*Result, error) {
	if conf.RuntimeID == "" && conf.ApplicationName == "" {
		logrus.Errorf("get-addon-info failed: need provide runtime_id or application_name.")
		return nil, errors.Errorf("get-addon-info failed: need provide runtime_id or application_name.")
	}

	if conf.AddonName == "" {
		logrus.Errorf("get-addon-info failed: runtime_id and addon_name must provided.")
		return nil, errors.Errorf("get-addon-info failed: runtime_id and addon_name must provided.")
	}

	result := &Result{}

	err := retry.DoWithInterval(func() error {
		result.Info = make(map[string]string)

		// runtime_id 未设置，根据 application_name 查找 runtime_id
		if conf.RuntimeID == "" {
			// 通过 ApplicationName 的方式仅支持 trantor 业务以及通 release 部署，走 pipeline 部署的还是需要提供 RuntimeID
			// 通过 application name 先获取到 Application ID，然后结合 WorkSpace 和  应用程序名称（通过 release 部署的） 获取到 runtimeID
			appId, err := getAppID(conf, conf.ApplicationName)
			if err != nil {
				logrus.Errorf("deploy failed: get app Id for appName %s failed, error: %v", conf.ApplicationName, err)
				return errors.Errorf("deploy failed: get app Id for appName %s failed, error: %v.", conf.ApplicationName, err)
			}

			runtimeId, err := getRuntimeId(conf, conf.ApplicationName, appId)
			if err != nil {
				logrus.Errorf("deploy failed: get runtime ID for appName %s failed, error: %v", conf.ApplicationName, err)
				return errors.Errorf("deploy failed: get runtime ID for appName %s failed, error: %v", conf.ApplicationName, err)
			}

			conf.RuntimeID = runtimeId
		}

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

func generateMetadata(result *Result) *metadata.Metadata {
	addrs := make([]metadata.MetadataField, 0)
	for name, addr := range result.Info {
		addrs = append(addrs, metadata.MetadataField{
			Name:  name,
			Value: addr,
		})
	}

	var ret metadata.Metadata
	ret = addrs
	return &ret
}

// 获取应用程序对应的 Runtime
func getAddons(conf *conf) ([]apistructs.AddonFetchResponseData, error) {
	var resp apistructs.AddonListResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPrefix).Path("/api/addons").
		Param("type", "runtime").
		Param("workspace", conf.Workspace).
		Param("value", conf.RuntimeID).
		Param("projectId", fmt.Sprintf("%d", conf.ProjectID)).
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
	if len(resp.Data) == 0 {
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
