package pkg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/command"
	"github.com/erda-project/erda-actions/pkg/dice"
	"github.com/erda-project/erda-actions/pkg/meta"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type Sonar struct {
	Auth    SonarAuth
	cfg     *Conf
	cmd     *command.Cmd
	results *meta.ResultMetaCollector
}

type SonarAuth struct {
	HostURL  string
	Login    string
	Password string
}

func NewSonar(cfg *Conf) (*Sonar, error) {
	if cfg.SonarHostURL == "" {
		app, err := dice.GetApplication(cfg.PlatformParams)
		if err != nil {
			return nil, err
		}
		if app.SonarConfig == nil {
			return nil, errors.Errorf("application %s has no sonar config", app.Name)
		}
		cfg.SonarHostURL = app.SonarConfig.Host
		cfg.SonarLogin = app.SonarConfig.Token
		cfg.ProjectKey = app.SonarConfig.ProjectKey
	}
	auth := SonarAuth{HostURL: cfg.SonarHostURL, Login: cfg.SonarLogin, Password: cfg.SonarPassword}
	sonar := Sonar{Auth: auth, cfg: cfg, cmd: command.NewCmd("sonar-scanner"), results: meta.NewResultMetaCollector()}
	return &sonar, nil
}

type SonarCeTaskStatus string

var (
	SonarCeTaskStatusSuccess    SonarCeTaskStatus = "SUCCESS"
	SonarCeTaskStatusFailed     SonarCeTaskStatus = "FAILED"
	SonarCeTaskStatusCanceled   SonarCeTaskStatus = "CANCELED"
	SonarCeTaskStatusInProgress SonarCeTaskStatus = "IN_PROGRESS"
)

// SonarCeTask sonar scanner task result
type SonarCeTask struct {
	AnalysisID      string            `json:"analysisId"`
	Status          SonarCeTaskStatus `json:"status"`
	Type            string            `json:"type"`
	ErrorMessage    string            `json:"errorMessage"`
	ErrorStacktrace string            `json:"errorStacktrace"`
}

// ScannerReportFileContent sonar scanner 分析时的报告文件，用于获取 ceTask
type ScannerReportFileContent struct {
	ProjectKey    string `json:"projectKey"`
	ServerURL     string `json:"serverUrl"`
	ServerVersion string `json:"serverVersion"`
	DashboardURL  string `json:"dashboardUrl"`
	CeTaskID      string `json:"ceTaskId"`
	CeTaskUrl     string `json:"ceTaskUrl"`
}

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

// handleScannerReportTaskFile 处理 report-task.txt 文件，使用 analysisID 获取 sonarCeTask 状态和错误信息
func (sonar *Sonar) handleScannerReportTaskFile(filePath string) (*ScannerReportFileContent, *SonarCeTask, error) {
	logrus.Infof("Begin handle scanner report task file...")
	defer func() {
		logrus.Infof("End handle scanner report task file")
	}()

	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open scanner report-task.txt, err: %v", err)
	}

	// 解析文件，k=v 格式，获取 report 内容
	logrus.Info("report-task.txt below:")
	fscan := bufio.NewScanner(f)
	var repoTask ScannerReportFileContent
	for fscan.Scan() {
		line := fscan.Text()
		fmt.Println(line)
		kv := strings.SplitN(line, "=", 2)
		if len(kv) < 2 {
			continue
		}
		switch kv[0] {
		case "projectKey":
			repoTask.ProjectKey = kv[1]
		case "serverUrl":
			repoTask.ServerURL = kv[1]
		case "serverVersion":
			repoTask.ServerVersion = kv[1]
		case "dashboardUrl":
			repoTask.DashboardURL = kv[1]
		case "ceTaskId":
			repoTask.CeTaskID = kv[1]
		case "ceTaskUrl":
			repoTask.CeTaskUrl = kv[1]
		}
	}
	logrus.Info("report-task.txt done")

	// 使用 ceTaskID 查询 ceTask
	var ceTask *SonarCeTask
	for {
		time.Sleep(time.Second * 5)
		ceTask, err = sonar.getSonarCeTaskResult(repoTask.CeTaskID)
		if err != nil {
			logrus.Infof("invoke sonar to get ce task status failed, continue, err: %v\n", err)
			continue
		}
		fmt.Println("ceTask below:")
		fmt.Println(jsonparse.JsonOneLine(ceTask))
		switch ceTask.Status {
		case SonarCeTaskStatusSuccess:
			logrus.Infof("ce task status: %s\n", ceTask.Status)
			return nil, ceTask, nil
		case SonarCeTaskStatusFailed, SonarCeTaskStatusCanceled:
			return nil, ceTask, fmt.Errorf("ce task status: %s, errMessage: %s, errStackTrace: %s\n", ceTask.Status, ceTask.ErrorMessage, ceTask.ErrorStacktrace)
		case SonarCeTaskStatusInProgress:
			logrus.Infof("ce task status: %s, waiting...\n", ceTask.Status)
		default:
			return nil, ceTask, fmt.Errorf("ce task status: %s, invalid", ceTask.Status)
		}
	}
}

func (sonar *Sonar) getSonarCeTaskResult(ceTaskID string) (*SonarCeTask, error) {
	hc := httpclient.New(httpclient.WithCompleteRedirect())
	var respBody bytes.Buffer
	resp, err := hc.BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).Path("/api/ce/task").Param("id", ceTaskID).
		Do().Body(&respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get ce task result, ceTaskID: %s, err: %v", ceTaskID, err)
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("failed to get ce task result, ceTaskID: %s, statusCode: %d, respBody: %s", ceTaskID, resp.StatusCode(), respBody.String())
	}
	var taskResp struct {
		Task SonarCeTask `json:"task"`
	}
	if err := json.Unmarshal(respBody.Bytes(), &taskResp); err != nil {
		return nil, fmt.Errorf("failed to parse ce task result, err: %v", err)
	}
	return &taskResp.Task, nil
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
