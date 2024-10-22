package pkg

import (
	"fmt"
	"path/filepath"

	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/sirupsen/logrus"
)

var (
	ResultKeyProjectKey        = "projectKey"
	ResultKeyLanguage          = "language"
	ResultKeyQualityGateStatus = "qualityGateStatus"
	ResultKeyQualityGateURL    = "qualityGateURL"
)

func (s *Sonar) Execute() error {
	if s.cfg.CodeDir == "" {
		return fmt.Errorf("missing context")
	}
	s.cmd.SetDir(s.cfg.CodeDir)
	s.cmd.Add(fmt.Sprintf("-Dsonar.host.url=%s", s.Auth.HostURL))
	s.cmd.Add(fmt.Sprintf("-Dsonar.login=%s", s.Auth.Login))
	if s.cfg.SonarPassword != "" {
		s.cmd.Add(fmt.Sprintf("-Dsonar.password=%s", s.cfg.SonarPassword))
	}

	// java
	if s.cfg.SonarJavaBinaries != "" {
		s.cmd.Add(fmt.Sprintf("-Dsonar.java.binaries=%s", s.cfg.SonarJavaBinaries))
	}

	projectKey := s.cfg.ProjectKey
	if projectKey == "" {
		return fmt.Errorf("missing project key")
	}
	logrus.Info("sonar project key: ", projectKey)
	s.cmd.Add(fmt.Sprintf("-Dsonar.projectKey=%s", projectKey))
	s.results.Add(ResultKeyProjectKey, projectKey)
	defer func() {
		if err := s.results.Store(); err != nil {
			logrus.Errorf("failed to store results: %v", err)
		}
		return
	}()
	s.cmd.Add("-Dsonar.scm.disabled=true")
	logrus.Infof("Begin analysis project by sonar-scanner...")
	defer func() {
		logrus.Infof("End analysis project by sonar-scanner")
	}()
	if err := s.cmd.Run(); err != nil {
		return fmt.Errorf("sonar analysis failed, err: %v", err)
	}

	logrus.Infof("Sonar Analysis finished, executing results...")
	_, ceTask, err := s.handleScannerReportTaskFile(filepath.Join(s.cfg.CodeDir, ".scannerwork/report-task.txt"))
	if err != nil {
		return fmt.Errorf("failed to analysis sonar project, err: %v", err)
	}
	s.results.Add(ResultKeyQualityGateURL, fmt.Sprintf("%s/dashboard?id=%s", s.cfg.SonarHostURL, s.cfg.ProjectKey))

	qualityGateResult, err := s.getSonarQualityGateStatus(ceTask.AnalysisID)
	if err != nil {
		return err
	}
	qualityGateStatus := qualityGateResult.Status
	s.results.Add(ResultKeyQualityGateStatus, string(qualityGateStatus))

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
		s.results.Add(fmt.Sprintf("metric: %s", condition.MetricKey), string(condition.Status))
		s.results.Add(fmt.Sprintf("metric(detail): %s", condition.MetricKey), jsonparse.JsonOneLine(condition))
	}

	if s.cfg.MustGateStatusOK && qualityGateStatus != SonarQualityGateStatusOK {
		logrus.Errorf("QUALITY GATE STATUS: %s", qualityGateStatus)
		return fmt.Errorf("sonar quality gate status is not ok, status: %s", qualityGateStatus)
	}
	return nil
}
