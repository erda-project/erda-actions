package push

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/api-publish/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

func Run() error {
	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		return errors.WithStack(err)
	}

	url := cfg.DiceOpenapiPrefix + "/api/gateway/registrations/" + cfg.RegisterId + "/publish"
	publishMsg := conf.PublishMsg{
		OrgId:      strconv.Itoa(int(cfg.OrgID)),
		OwnerEmail: cfg.OwnerEmail,
		ItemName:   cfg.ItemName,
	}
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
