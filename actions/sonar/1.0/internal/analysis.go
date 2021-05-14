package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/jsonparse"
	"github.com/erda-project/erda/pkg/uuid"
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
		projectKey = uuid.UUID()
	}
	logrus.Infof("sonar.projectKey: %s", projectKey)
	args = append(args, fmt.Sprintf("-Dsonar.projectKey=%s", projectKey))
	// defer delete project after analyze
	if cfg.ActionParams.DeleteProject {
		defer sonar.invokeDelSonarServerProject(projectKey)
	}
	results.Add(ResultKeyProjectKey, projectKey)

	// language
	if !cfg.ActionParams.Language.Supported() {
		return nil, fmt.Errorf("sonar analysis not supported language %s yet", cfg.ActionParams.Language)
	}
	switch cfg.ActionParams.Language {
	case LanguageGo:
	case LanguageJava:
		// sonar.java.binaries
		sonarJavaBinary := cfg.ActionParams.SonarJavaBinaries
		if sonarJavaBinary == "" {
			sonarJavaBinary = "."
		}
		args = append(args, fmt.Sprintf("-Dsonar.java.binaries=%s", sonarJavaBinary))
	}
	results.Add(ResultKeyLanguage, string(cfg.ActionParams.Language))

	// exclusions
	if cfg.ActionParams.SonarExclusions != "" {
		args = append(args, fmt.Sprintf("-Dsonar.exclustions=%s", cfg.ActionParams.SonarExclusions))
	}

	// scm
	args = append(args, "-Dsonar.scm.disabled=true")

	// log
	args = append(args, fmt.Sprintf("-Dsonar.log.level=%s", cfg.ActionParams.SonarLogLevel))
	switch cfg.ActionParams.SonarLogLevel {
	case SonarLogLevelDEBUG, SonarLogLevelTRACE:
		args = append(args, fmt.Sprintf("-Dsonar.verbose=%t", true))
	}

	// TODO
	// produceCoverageArgs(lan, cfg.ActionParams.GoDir, contextDir)

	// 在分析前创建项目，因为 sonar-scanner cli 无法指定 自定义质量门禁，只能先创建项目和自定义门禁后，再进行绑定
	if err := sonar.createProject(projectKey); err != nil {
		return nil, fmt.Errorf("failed to create sonar project for analysis, projectKey: %s, err: %v", projectKey, err)
	}

	// 获取平台级配置
	var keys []*apistructs.SonarMetricKey
	var err error
	if cfg.ActionParams.UsePlatformQualityGate {
		// 获取平台级配置
		keys, err = sonar.getSonarMetricKeys(cfg)
		if err != nil {
			return nil, fmt.Errorf("get platform metric key value error %v", err)
		}

		fmt.Println("load platform metric keys success")
		spew.Dump(keys)
	}

	// 假如用户填了就不把平台的添加进去
	for _, v := range keys {
		var find = false
		for _, gate := range cfg.ActionParams.QualityGate {
			if v.MetricKey == gate.Metric.String() {
				find = true
				break
			}
		}
		// 用户没填的就加入进去
		if !find {
			cfg.ActionParams.QualityGate = append(cfg.ActionParams.QualityGate, QualityGateCondition{
				Metric: MetricKey(v.MetricKey),
				Op:     QualityGateConditionOp(v.Operational),
				Error:  v.MetricValue,
			})
		}
	}

	// custom quality gate
	if len(cfg.ActionParams.QualityGate) > 0 {

		gateName := projectKey
		// create qg
		if err := sonar.createQualityGate(gateName); err != nil {
			return nil, fmt.Errorf("failed to create sonar custom quality gate, err: %v", err)
		}
		if cfg.ActionParams.DeleteProject {
			defer func() {
				if err := sonar.destroyQualityGate(gateName); err != nil {
					logrus.Warnf("failed to create sonar custom quality gate, err: %v", err)
				}
			}()
		}
		// create qg conditions
		for _, cond := range cfg.ActionParams.QualityGate {
			if err := sonar.createQualityGateCondition(gateName, cond); err != nil {
				return nil, fmt.Errorf("failed to create sonar custom quality gate condition, metric: %s, err: %v", cond.Metric.String(), err)
			}
		}
		// associate project to qg
		if err := sonar.associateProjectToQualityGate(projectKey, gateName); err != nil {
			return nil, fmt.Errorf("failed to associate project to custom quality gate, projectKey: %s, gateName: %s, err: %v", projectKey, gateName, err)
		}
	}

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

	// report sonar issues to qa
	if _, err := sonar.reportSonarIssues2QA(projectKey, cfg); err != nil {
		logrus.Infof("Failed to report sonar issues to QA Platform, err: %v", err)
		return nil, fmt.Errorf("failed to report sonar issue to QA Platform, projectKey: %s, err: %v", projectKey, err)
	}

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
	if qualityGateStatus != SonarQualityGateStatusOK {
		logrus.Errorf("QUALITY GATE STATUS: %s", qualityGateStatus)
		return results, fmt.Errorf("quality gate status: %s", qualityGateStatus)
	}

	return results, nil
}
