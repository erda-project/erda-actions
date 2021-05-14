package main

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/httpclient"
)

type SonarQualityGateResult struct {
	Status     SonarQualityGateStatus            `json:"status"`
	Conditions []SonarQualityGateConditionResult `json:"conditions"`
}
type SonarQualityGateConditionResult struct {
	Status         SonarQualityGateStatus `json:"status"`
	MetricKey      MetricKey              `json:"metricKey"`
	Comparator     string                 `json:"comparator"`
	ErrorThreshold string                 `json:"errorThreshold"`
	ActualValue    string                 `json:"actualValue"`
}

type SonarQualityGateStatus string

var (
	SonarQualityGateStatusOK    SonarQualityGateStatus = "OK"
	SonarQualityGateStatusWARN  SonarQualityGateStatus = "WARN"
	SonarQualityGateStatusERROR SonarQualityGateStatus = "ERROR"
	SonarQualityGateStatusNONE  SonarQualityGateStatus = "NONE"
)

// getSonarQualityGateStatus 查询代码门禁结果
func (sonar *Sonar) getSonarQualityGateStatus(analysisID string) (*SonarQualityGateResult, error) {
	var statusResp struct {
		ProjectStatus *SonarQualityGateResult `json:"projectStatus"`
	}
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).Path("/api/qualitygates/project_status").
		Param("analysisId", analysisID).
		Do().JSON(&statusResp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("failed to fetch qualitygates status, status code: %d", resp.StatusCode())
	}
	return statusResp.ProjectStatus, nil
}

// createQualityGate 创建代码质量门禁
func (sonar *Sonar) createQualityGate(name string) error {
	logrus.Infof("Begin create sonar quality gate: %s", name)
	var respBody bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Post(sonar.Auth.HostURL).Path("/api/qualitygates/create").
		Param("name", name).
		Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !httpResp.IsOK() {
		logrus.Errorf("Failed to create sonar quality gate: %s, err: %v", name, err)
		return fmt.Errorf("failed to create quality gate, name: %s, status-code: %d, respBody: %s", name, httpResp.StatusCode(), respBody.String())
	}
	logrus.Infof("End create sonar quality gate: %s", name)
	return nil
}

// destroyQualityGate 销毁代码质量门禁
func (sonar *Sonar) destroyQualityGate(name string) error {
	var respBody bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Post(sonar.Auth.HostURL).Path("/api/qualitygates/destroy").
		Param("name", name).
		Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !httpResp.IsOK() {
		return fmt.Errorf("failed to destroy quality gate, name: %s, status-code: %d, respBody: %s", name, httpResp.StatusCode(), respBody.String())
	}
	return nil
}

type QualityGateConditionOp string

func (op QualityGateConditionOp) String() string {
	return string(op)
}

var (
	OpGT QualityGateConditionOp = "GT" // is greater than
	OpLT QualityGateConditionOp = "LT" // is lower than
)

type QualityGateCondition struct {
	Error  string                 `json:"error"`  // condition error threshold. "A","B" will transfer to 1,2
	Metric MetricKey              `json:"metric"` // condition metric
	Op     QualityGateConditionOp `json:"op"`     // condition operator
}

// createQualityGateCondition
func (sonar *Sonar) createQualityGateCondition(gateName string, cond QualityGateCondition) error {
	logrus.Infof("Begin create sonar quality gate condition, metric: %s, op: %s, err: %s", cond.Metric.String(), cond.Op, cond.Error)
	// generate params
	params := make(url.Values, 0)
	params.Add("gateName", gateName)
	params.Add("error", cond.Error)
	params.Add("metric", cond.Metric.String())
	params.Add("op", cond.Op.String())

	var respBody bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Post(sonar.Auth.HostURL).Path("/api/qualitygates/create_condition").
		Params(params).
		Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !httpResp.IsOK() {
		logrus.Errorf("Failed to create sonar quality gate condition, metric: %s, err: %v", cond.Metric, err)
		return fmt.Errorf("failed to create quality gate condition, gateName: %s, errorThreshold: %s, metric: %s, op: %s, respBody: %s",
			gateName, cond.Error, cond.Metric, cond.Op, respBody.String())
	}
	logrus.Infof("End create sonar quality gate condition, metric: %s", cond.Metric)
	return nil
}
