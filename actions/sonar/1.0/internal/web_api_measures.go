package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/erda-project/erda/apistructs"

	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
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

// ProjectMeasures 项目测量值
type ProjectMeasures struct {
	Component struct {
		Key      string `json:"key"`
		Measures []struct {
			Metric    MetricKey `json:"metric"`
			Value     string    `json:"value"`
			BestValue bool      `json:"bestValue"`
		} `json:"measures"`
	} `json:"component"`
}

// querySonarProjectMeasures 查询 sonar 项目测量值
// web api:/api/measures/component
func (sonar *Sonar) querySonarProjectMeasures(projectKey string) (*ProjectMeasures, error) {
	metricKeys := []MetricKey{
		MetricKeyBugs, MetricKeyVulnerabilities, MetricKeyCodeSmells,
		MetricKeyCoverage,
		MetricKeyNcloc, MetricKeyDuplicatedLinesDensity, MetricKeyDuplicatedBlocks,
		MetricKeyReliabilityRating,
		MetricKeySecurityRating,
		MetricKeyMaintainabilityRating,
	}
	var metricKeyStrs []string
	for _, key := range metricKeys {
		metricKeyStrs = append(metricKeyStrs, string(key))
	}
	var body bytes.Buffer
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).
		Path("/api/measures/component").
		Param("additionalFields", "metrics").
		Param("component", projectKey).
		Param("metricKeys", strutil.Join(metricKeyStrs, ",", true)).
		Do().
		Body(&body)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("statusCode: %d, responseBody: %s", resp.StatusCode(), body.String())
	}

	var measures ProjectMeasures
	if err := json.Unmarshal(body.Bytes(), &measures); err != nil {
		return nil, err
	}

	return &measures, nil
}

func makeQATestIssueStatistics(pm ProjectMeasures) apistructs.TestIssuesStatistics {
	statistic := apistructs.TestIssuesStatistics{
		Rating: &apistructs.TestIssueStatisticsRating{},
	}
	for _, measure := range pm.Component.Measures {
		switch measure.Metric {
		case MetricKeyBugs:
			statistic.Bugs = measure.Value
		case MetricKeyVulnerabilities:
			statistic.Vulnerabilities = measure.Value
		case MetricKeyCodeSmells:
			statistic.CodeSmells = measure.Value
		case MetricKeyCoverage:
			statistic.Coverage = measure.Value
		case MetricKeyDuplicatedLinesDensity:
			statistic.Duplications = measure.Value
		case MetricKeyReliabilityRating:
			statistic.Rating.Bugs = ratingConvert(measure.Value)
		case MetricKeySecurityRating:
			statistic.Rating.Vulnerabilities = ratingConvert(measure.Value)
		case MetricKeyMaintainabilityRating:
			statistic.Rating.CodeSmells = ratingConvert(measure.Value)
		}
	}
	statistic.SonarKey = pm.Component.Key
	return statistic
}

func ratingConvert(rating string) apistructs.CodeQualityRatingLevel {
	n, _ := strconv.ParseFloat(rating, 10)
	fmt.Printf("rating convert, rating: %s, parsed: %v\n", rating, n)
	switch n {
	case 1.0:
		return apistructs.CodeQualityRatingLevelA
	case 2.0:
		return apistructs.CodeQualityRatingLevelB
	case 3.0:
		return apistructs.CodeQualityRatingLevelC
	case 4.0:
		return apistructs.CodeQualityRatingLevelD
	case 5.0:
		return apistructs.CodeQualityRatingLevelE
	default:
		return apistructs.CodeQualityRatingLevelUnknown
	}
}
