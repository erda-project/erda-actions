package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

// reportSonarIssues2QA 将 sonar issue 转换为 qa issue 并上报
func (sonar *Sonar) reportSonarIssues2QA(projectKey string, cfg *Conf) (*apistructs.SonarStoreRequest, error) {
	logrus.Infof("Begin report sonar issues to QA Platform...")
	defer func() {
		logrus.Infof("End report sonar issue to QA Platform")
	}()

	var store apistructs.SonarStoreRequest
	store.Key = projectKey
	store.CommitID = cfg.GittarCommit
	store.LogID = cfg.LogID
	store.Branch = cfg.GittarBranch
	store.GitRepo = cfg.GittarRepo
	store.ApplicationID = int64(cfg.AppID)
	store.ApplicationName = cfg.AppName
	store.BuildID = cfg.PipelineID
	store.ProjectID = int64(cfg.ProjectID)
	store.ProjectName = cfg.ProjectName
	store.OperatorID = cfg.OperatorID

	// 统计指标
	projectMeasures, err := sonar.querySonarProjectMeasures(projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to query soanr project measures, err: %v", err)
	}
	statistic := makeQATestIssueStatistics(*projectMeasures)
	store.IssuesStatistics = statistic

	// sonar issues 可能较多，分 issue type 多次查询
	// get bugs
	bugs, err := sonar.convert2QaIssues(projectKey, IssueTypeBug)
	if err != nil {
		return nil, err
	}
	store.Bugs = bugs

	// get vulnerabilities
	vulnerabilities, err := sonar.convert2QaIssues(projectKey, IssueTypeVulnerability)
	if err != nil {
		return nil, err
	}
	store.Vulnerabilities = vulnerabilities

	// get codeSmells
	codeSmells, err := sonar.convert2QaIssues(projectKey, IssueTypeCodeSmell)
	if err != nil {
		return nil, err
	}
	store.CodeSmells = codeSmells

	// get coverage
	coverage, err := sonar.getQaIssueTree(projectKey, MeasureTypeCoverage)
	if err != nil {
		return nil, err
	}
	store.Coverage = coverage

	// get duplications
	duplications, err := sonar.getQaIssueTree(projectKey, MeasureTypeDuplications)
	if err != nil {
		return nil, err
	}
	store.Duplications = duplications

	if err := retry.DoWithInterval(func() error {
		var body bytes.Buffer
		r, err := httpclient.New(httpclient.WithCompleteRedirect(), httpclient.WithTimeout(30*time.Second, 180*time.Second)).
			Post(cfg.OpenAPIAddr).
			Path("/api/qa/actions/sonar-results-store").
			JSONBody(store).
			Header("Content-Type", "application/json").
			Header("Authorization", cfg.OpenAPIToken).
			Do().
			Body(&body)
		if err != nil {
			return err
		}
		if !r.IsOK() {
			return errors.Errorf("statusCode: %d, responseBody: %s", r.StatusCode(), body.String())
		}
		return nil
	}, 2, time.Second*5); err != nil {
		return nil, fmt.Errorf("report sonar issues to platform failed, projectKey: %s, err: %v", projectKey, err)
	} else {
		logrus.Infof("report sonar issues to platform success, projectKey: %s", projectKey)
	}

	return &store, nil
}

// convert2QaIssues 将 sonar issues 转换为 qa issues
func (sonar *Sonar) convert2QaIssues(projectKey string, issueType IssueType) ([]*apistructs.TestIssues, error) {

	// sonar issues
	sonarIssues, err := sonar.querySonarIssues(projectKey, issueType)
	if err != nil {
		return nil, err
	}

	// sonar issues -> qa issues
	var qaIssues []*apistructs.TestIssues
	for _, sonarIssue := range sonarIssues {
		var codes []string
		for _, sourceLine := range sonarIssue.SourceLines {
			codes = append(codes, sourceLine.Code)
		}
		qaIssue := apistructs.TestIssues{
			Path:      sonarIssue.ComponentPath,
			Component: sonarIssue.Component,
			Message:   sonarIssue.Message,
			Rule:      sonarIssue.Rule,
			TextRange: sonarIssue.TextRange,
			Severity:  sonarIssue.Severity.String(),
			Line:      sonarIssue.LineNum,
			Code:      codes,
		}
		qaIssues = append(qaIssues, &qaIssue)
	}

	return qaIssues, nil
}
