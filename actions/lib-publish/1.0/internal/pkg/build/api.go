package build

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/lib-publish/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

// GetAppPublishItemRelations 根据 appID 获取应用关联的发布内容
func GetAppPublishItemRelations(cfg conf.Conf) (*apistructs.QueryAppPublishItemRelationResponse, error) {
	var resp apistructs.QueryAppPublishItemRelationResponse
	request := httpclient.New(httpclient.WithCompleteRedirect()).Get(cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/applications/%d/actions/get-publish-item-relations", cfg.AppID)).
		Header("Authorization", cfg.CiOpenapiToken)
	httpResp, err := request.Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, errors.Errorf("failed to get publish item relation, status code: %d, body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return nil, errors.Errorf(resp.Error.Msg)
	}
	return &resp, nil
}

// CreatePublishItemVersion 创建发布内容版本
func CreatePublishItemVersion(cfg conf.Conf, req apistructs.CreatePublishItemVersionRequest) (*apistructs.CreatePublishItemVersionResponse, error) {
	var resp apistructs.CreatePublishItemVersionResponse
	request := httpclient.New(httpclient.WithCompleteRedirect()).Post(cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/publish-items/%d/versions", req.PublishItemID)).
		Header("Authorization", cfg.CiOpenapiToken)
	httpResp, err := request.JSONBody(&req).Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, errors.Errorf("failed to create publish item version, status code: %d, body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return nil, errors.Errorf(resp.Error.Msg)
	}
	return &resp, nil
}
