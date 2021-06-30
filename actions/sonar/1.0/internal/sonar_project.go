package main

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/http/httpclient"
)

// createProject 创建 sonar 项目
func (sonar *Sonar) createProject(projectKey string) error {
	logrus.Infof("Begin create sonar project, projectKey: %s", projectKey)
	var respBody bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Post(sonar.Auth.HostURL).Path("/api/projects/create").
		Param("name", projectKey).
		Param("project", projectKey).
		Param("visibility", "private").
		Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !httpResp.IsOK() {
		logrus.Errorf("Failed to create sonar project, err: %v", err)
		return fmt.Errorf("failed to create project, projectKey: %s, status-code: %d, respBody: %s", projectKey, httpResp.StatusCode(), respBody.String())
	}
	logrus.Infof("End create sonar project")
	return nil
}

// associateProjectToQualityGate Associate a project to a quality gate.
func (sonar *Sonar) associateProjectToQualityGate(projectKey string, gateName string) error {
	logrus.Infof("Begin associate sonar project %s to quality gate %s", projectKey, gateName)
	defer func() {
		logrus.Errorf("End associate sonar project to quality gate")
	}()
	var respBody bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Post(sonar.Auth.HostURL).Path("/api/qualitygates/select").
		Param("gateName", gateName).
		Param("projectKey", projectKey).
		Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !httpResp.IsOK() {
		return fmt.Errorf("failed to associate project to quality gate, projectKey: %s, gateName: %s, status-code: %d, respBody: %s", projectKey, gateName, httpResp.StatusCode(), respBody.String())
	}
	return nil
}
