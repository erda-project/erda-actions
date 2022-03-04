package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/encoding/jsonparse"
)

// Analysis 使用 sonar-scanner 进行代码质量分析
func (sonar *Sonar) Analysis(cfg *Conf) (*ResultMetas, error) {
	// 存储结果
	results := &ResultMetas{}

	// make sonar-scanner args
	scanner := exec.Command("sonar-scanner")
	var args []string

	// exec context
	if cfg.ActionParams.CodeDir == "" {
		return nil, fmt.Errorf("missing context")
	}
	scanner.Dir = cfg.ActionParams.CodeDir

	// sonar server
	args = append(args, fmt.Sprintf("-Dsonar.host.url=%s", sonar.Auth.HostURL))

	// sonar authentication
	args = append(args, fmt.Sprintf("-Dsonar.login=%s", sonar.Auth.Login))
	if cfg.ActionParams.SonarPassword != "" {
		args = append(args, fmt.Sprintf("-Dsonar.password=%s", sonar.Auth.Password))
	}

	// project key
	projectKey := cfg.ActionParams.ProjectKey
	if projectKey == "" {
		return nil, fmt.Errorf("missing project key")
	}
	logrus.Infof("sonar.projectKey: %s", projectKey)
	args = append(args, fmt.Sprintf("-Dsonar.projectKey=%s", projectKey))
	results.Add(ResultKeyProjectKey, projectKey)

	// scm
	args = append(args, "-Dsonar.scm.disabled=true")

	// begin sonar-scanner
	logrus.Infof("Begin analysis project by sonar-scanner...")
	defer func() {
		logrus.Infof("End analysis project by sonar-scanner")
	}()
	scanner.Args = append(scanner.Args, args...)
	if cfg.Debug {
		fmt.Println(scanner.String())
	}
	scanner.Stdout = os.Stdout
	scanner.Stderr = os.Stderr
	if err := scanner.Run(); err != nil {
		return nil, fmt.Errorf("sonar analysis failed, err: %v", err)
	}

	logrus.Infof("Sonar Analysis finished, 正在处理结果...")
	// handle scanner report file
	reportFile, ceTask, err := sonar.handleScannerReportTaskFile(filepath.Join(cfg.ActionParams.CodeDir, ".scannerwork/report-task.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to analysis sonar project, err: %v", err)
	}
	_ = reportFile

	// quality gate url
	results.Add(ResultKeyQualityGateURL, fmt.Sprintf("%s/dashboard?id=%s", cfg.ActionParams.SonarHostURL, cfg.ActionParams.ProjectKey))
	// quality gate
	qualityGateResult, err := sonar.getSonarQualityGateStatus(ceTask.AnalysisID)
	if err != nil {
		return results, err
	}
	qualityGateStatus := qualityGateResult.Status
	results.Add(ResultKeyQualityGateStatus, string(qualityGateStatus))
	// add condition results to meta
	fmt.Println("quality gate conditions result below:")
	for _, condition := range qualityGateResult.Conditions {
		// log
		fmt.Printf("metric: %s, status: %s\n", condition.MetricKey, condition.Status)
		fmt.Printf("metric(detail): %s, detail: %s\n", condition.MetricKey, jsonparse.JsonOneLine(condition))
		// if not ok, add to meta and show it
		if condition.Status == SonarQualityGateStatusOK {
			continue
		}
		results.Add(ResultKey(fmt.Sprintf("metric: %s", condition.MetricKey)), string(condition.Status))
		results.Add(ResultKey(fmt.Sprintf("metric(detail): %s", condition.MetricKey)), jsonparse.JsonOneLine(condition))
	}

	if cfg.ActionParams.MustGateStatusOK && qualityGateStatus != SonarQualityGateStatusOK {
		logrus.Errorf("QUALITY GATE STATUS: %s", qualityGateStatus)
		return results, fmt.Errorf("sonar quality gate status is not ok, status: %s", qualityGateStatus)
	}

	return results, nil
}
