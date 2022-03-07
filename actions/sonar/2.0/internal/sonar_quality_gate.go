package main

import (
	"fmt"

	"github.com/erda-project/erda/pkg/http/httpclient"
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

// MetricKey 指标名
// see: https://docs.sonarqube.org/latest/user-guide/metric-definitions/
type MetricKey string

var (
	MetricKeyBugs                   MetricKey = "bugs"
	MetricKeyVulnerabilities        MetricKey = "vulnerabilities"
	MetricKeyCodeSmells             MetricKey = "code_smells"
	MetricKeyCoverage               MetricKey = "coverage"
	MetricKeyUncoveredLines         MetricKey = "uncovered_lines"
	MetricKeyUncoveredConditions    MetricKey = "uncovered_conditions"
	MetricKeyNcloc                  MetricKey = "ncloc"                    // Lines of code: Number of physical lines that contain at least one character which is neither a whitespace nor a tabulation nor part of a comment.
	MetricKeyDuplicatedLinesDensity MetricKey = "duplicated_lines_density" // 重复行数百分比：duplicated_lines / lines * 100
	MetricKeyDuplicatedLines        MetricKey = "duplicated_lines"
	MetricKeyDuplicatedBlocks       MetricKey = "duplicated_blocks"
	MetricKeyReliabilityRating      MetricKey = "reliability_rating"
	MetricKeySecurityRating         MetricKey = "security_rating"
	MetricKeyMaintainabilityRating  MetricKey = "sqale_rating"

	/////////////////////////////////////////////
	// quality gate built-in sonar-way metrics //
	/////////////////////////////////////////////
	MetricKeyNewSecurityRating           MetricKey = "new_security_rating"
	MetricKeyNewReliabilityRating        MetricKey = "new_reliability_rating"
	MetricKeyNewMaintainabilityRating    MetricKey = "new_maintainability_rating"
	MetricKeyNewCoverage                 MetricKey = "new_coverage"
	MetricKeyNewDuplicatedLinesDensity   MetricKey = "new_duplicated_lines_density"
	MetircKeyNewSecurityHotspotsReviewed MetricKey = "new_security_hotspots_reviewed"
)

func (k MetricKey) String() string {
	return string(k)
}

// MeasureType 测量类型
type MeasureType string

var (
	MeasureTypeCoverage     MeasureType = "coverage"
	MeasureTypeDuplications MeasureType = "duplications"
)

func (t MeasureType) String() string {
	return string(t)
}

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
