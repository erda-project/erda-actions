package push

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/erda-project/erda-actions/actions/api-register/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"

	"github.com/pkg/errors"
)

func Run() error {
	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		return errors.WithStack(err)
	}
	sjson, err := ioutil.ReadFile(cfg.SwaggerPath)
	if err != nil {
		return errors.WithStack(err)
	}
	var swaggerJson interface{}
	err = json.Unmarshal(sjson, &swaggerJson)
	if err != nil {
		return errors.WithStack(err)
	}
	registerMsg := conf.RegisterApiMsg{
		OrgId:       strconv.FormatInt(cfg.OrgID, 10),
		ProjectId:   strconv.FormatInt(cfg.ProjectID, 10),
		Workspace:   cfg.Workspace,
		ClusterName: cfg.ClusterName,
		AppId:       strconv.FormatInt(cfg.AppID, 10),
		AppName:     cfg.AppName,
		RuntimeId:   cfg.RuntimeID,
		RuntimeName: cfg.GittarBranch,
		ServiceName: cfg.ServiceName,
		ServiceAddr: cfg.ServiceAddr,
		Swagger:     swaggerJson,
	}
	body, err := json.Marshal(registerMsg)
	if err != nil {
		return errors.WithStack(err)
	}
	url := cfg.DiceOpenapiPrefix + "/api/gateway/registrations"
	headers := make(map[string]string)
	headers["Authorization"] = cfg.CiOpenapiToken
	resp, _, err := Request("POST", url, body, 60, headers)
	if err != nil {
		return err
	}
	result := conf.HttpResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return errors.Wrapf(err, "resp:%s", resp)
	}
	if !result.Success {
		return errors.Errorf("register failed, resp:%s", resp)
	}
	for {
		complete, err := GetPublishStatus(cfg, result.Data.ApiRegisterId)
		if err != nil {
			return err
		}
		if complete {
			fmt.Fprintln(os.Stdout, "api register done")
			metaContent := fmt.Sprintf("registerId=%s\n", result.Data.ApiRegisterId)
			err = filehelper.CreateFile(cfg.Metafile, metaContent, 0644)
			if err != nil {
				return errors.WithStack(err)
			}
			break
		}
		fmt.Fprintln(os.Stdout, "api register not ready")
		time.Sleep(5 * time.Second)
	}
	return nil
}
