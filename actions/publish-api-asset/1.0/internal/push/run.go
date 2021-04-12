package push

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"

	"github.com/erda-project/erda-actions/actions/publish-api-asset/1.0/internal/conf"
)

func versionExtract(version string) (major, minor, patch uint64, err error) {
	items := strings.Split(version, ".")
	if len(items) != 3 {
		goto failed
	}
	major, err = strconv.ParseUint(items[0], 10, 64)
	if err != nil {
		goto failed
	}
	minor, err = strconv.ParseUint(items[1], 10, 64)
	if err != nil {
		goto failed
	}
	patch, err = strconv.ParseUint(items[2], 10, 64)
	if err != nil {
		goto failed
	}
	return
failed:
	err = errors.Errorf("Error publishing to Exchange: Invalid version '%s', expecting 'x.y.z' format,"+
		" Examples of good versions are 1.0.0 or 4.3.1", version)
	return
}

func Run() error {
	var cfg conf.Conf
	var err error
	if err = envconf.Load(&cfg); err != nil {
		return errors.WithStack(err)
	}
	url := cfg.DiceOpenapiPrefix + "/api/api-assets"
	var major, minor, patch uint64
	if cfg.Version != "" {
		major, minor, patch, err = versionExtract(cfg.Version)
		if err != nil {
			return errors.Cause(err)
		}
	}
	versionInfo := apistructs.APIAssetVersionCreateRequest{
		Major: major,
		Minor: minor,
		Patch: patch,
		Instances: []apistructs.APIAssetVersionInstanceCreateRequest{
			{
				InstanceType: apistructs.APIInstanceTypeService,
				RuntimeID:    cfg.RuntimeID,
				ServiceName:  cfg.ServiceName,
			},
		},
	}
	publishMsg := apistructs.APIAssetCreateRequest{
		OrgID:     cfg.OrgID,
		ProjectID: cfg.ProjectID,
		AppID:     cfg.AppID,
		AssetID:   apistructs.APIAssetID(cfg.AssetID),
		AssetName: cfg.DisplayName,
		Source:    "action",
		IdentityInfo: apistructs.IdentityInfo{
			UserID: cfg.UserID,
		},
	}
	if publishMsg.AssetName == "" {
		publishMsg.AssetName = cfg.AssetID
	}
	specContent, err := ioutil.ReadFile(cfg.SpecPath)
	if err != nil {
		return errors.WithStack(err)
	}
	versionInfo.Spec = string(specContent)
	if strings.HasSuffix(cfg.SpecPath, ".yaml") || strings.HasSuffix(cfg.SpecPath, ".yml") {
		versionInfo.SpecProtocol = apistructs.APISpecProtocolOas2Yaml
	} else if strings.HasSuffix(cfg.SpecPath, ".json") || strings.HasSuffix(cfg.SpecPath, ".js") {
		versionInfo.SpecProtocol = apistructs.APISpecProtocolOas2Json
	} else {
		return errors.Errorf("unknown api specification protocol, the file's name need end in .yaml or .json, file: %s ", cfg.SpecPath)
	}
	publishMsg.Versions = []apistructs.APIAssetVersionCreateRequest{versionInfo}
	body, err := json.Marshal(publishMsg)
	if err != nil {
		return errors.Wrapf(err, "body:%s", body)
	}
	headers := make(map[string]string)
	headers["Authorization"] = cfg.CiOpenapiToken
	result, _, err := Request("POST", url, body, 60, headers)
	if err != nil {
		return err
	}
	resp := conf.HttpResponse{}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return errors.Wrapf(err, "resp:%s", result)
	}
	if !resp.Success {
		return errors.Errorf("error response:%s", result)
	}
	return nil
}
