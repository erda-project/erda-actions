package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
)

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
