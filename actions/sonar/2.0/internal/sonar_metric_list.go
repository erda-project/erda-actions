package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

// 查询门禁
func (sonar *Sonar) getSonarMetricKeys(cfg *Conf) ([]*apistructs.SonarMetricKey, error) {

	projectID := strconv.Itoa(int(cfg.ProjectID))

	var buffer bytes.Buffer
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(cfg.OpenAPIAddr).
		Path("/api/sonar-metric-rules/actions/query-list").
		Param("scopeType", "project").
		Param("scopeId", projectID).
		Header("Authorization", cfg.OpenAPIToken).
		Do().Do().Body(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to getSonarMetricKeys, err: %v", err)
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("failed to getSonarMetricKeys, statusCode: %d, respBody: %s", resp.StatusCode(), buffer.String())
	}
	var result apistructs.SonarMetricRulesListResp
	respBody := buffer.String()
	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return nil, fmt.Errorf("failed to getSonarMetricKeys, err: %v, json string: %s", err, respBody)
	}

	return result.Results, nil
}
