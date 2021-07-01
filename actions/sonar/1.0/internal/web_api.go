// this file contains all sonar web api invokes
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

type Components struct {
	Key      string                     `json:"key"`  // projectKey:filePath
	Path     string                     `json:"path"` // filePath
	Name     string                     `json:"name"` // fileBaseName
	Language Language                   `json:"language"`
	Measures []*apistructs.TestMeasures `json:"measures"`
}

type SonarMeasuresComponentTree struct {
	Components []Components `json:"components"`
}

// invokeSonarIssuesTree
// web api: /api/measures/component_tree
func (sonar *Sonar) invokeSonarIssuesTree(projectKey string, measureType MeasureType) (*SonarMeasuresComponentTree, error) {
	var (
		metricKeys string
		metricSort string
	)
	switch measureType {
	case MeasureTypeCoverage:
		metricKeys = strings.Join([]string{
			MetricKeyCoverage.String(),
			MetricKeyUncoveredLines.String(),
			MetricKeyUncoveredConditions.String(),
		}, ",")
		metricSort = MetricKeyUncoveredLines.String()
	case MeasureTypeDuplications:
		metricKeys = strings.Join([]string{
			MetricKeyDuplicatedLinesDensity.String(),
			MetricKeyDuplicatedLines.String(),
		}, ",")
		metricSort = MetricKeyDuplicatedLinesDensity.String()
	default:
		return nil, errors.Errorf("error type, key: %s", metricKeys)
	}

	var body bytes.Buffer
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).
		Path("/api/measures/component_tree").
		Param("baseComponentKey", projectKey).
		Param("metricSortFilter", "withMeasuresOnly").
		Param("ps", "500").
		Param("asc", "false"). // desc
		Param("metricSort", metricSort).
		Param("s", "metric").
		Param("metricKeys", metricKeys).
		Param("strategy", "leaves").
		Do().
		Body(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to get sonar issue tree, projectKey: %s, type: %s, err: %v", projectKey, measureType, err)
	}
	if !r.IsOK() {
		return nil, fmt.Errorf("failed to get sonar issue tree, projectKey: %s, type: %s, statusCode: %d, responseBody: %s", projectKey, measureType, r.StatusCode(), body.String())
	}
	var tree SonarMeasuresComponentTree
	if err := json.Unmarshal(body.Bytes(), &tree); err != nil {
		return nil, fmt.Errorf("failed to parse sonar issue tree, projectKey: %s, type: %s, err: %v", projectKey, measureType, err)
	}

	return &tree, nil
}

// invokeDelSonarServerProject 删除用于分析的 sonar project
// web api: /api/projects/delete
func (sonar *Sonar) invokeDelSonarServerProject(projectKey string) error {
	req := make(url.Values)
	req.Set("project", projectKey)

	err := retry.DoWithInterval(func() error {
		var body bytes.Buffer
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).
			BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
			Post(sonar.Auth.HostURL).
			Path("/api/projects/delete").
			FormBody(req).
			// Header("Authentication", "admin").
			Do().
			Body(&body)
		if err != nil {
			return fmt.Errorf("delete sonar project failed, err: %v, url: %s, body: %s", err, sonar.Auth.HostURL, body.String())
		}
		if !r.IsOK() {
			return fmt.Errorf("statusCode: %d, body: %s", r.StatusCode(), body.String())
		}

		return nil
	}, 2, time.Second*2)
	if err != nil {
		return fmt.Errorf("failed to delete sonar project, projectKey: %s, err: %v", projectKey, err)
	}
	logrus.Infof("delete sonar project: %s successfully!", projectKey)
	return nil
}
