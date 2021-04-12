package main

import (
	"strconv"

	"github.com/erda-project/erda/apistructs"
)

type Sonar struct {
	Auth SonarAuth
}

type SonarAuth struct {
	HostURL  string
	Login    string
	Password string
}

func NewSonar(hostURL, login, password string) *Sonar {
	auth := SonarAuth{HostURL: hostURL, Login: login, Password: password}
	sonar := Sonar{Auth: auth}
	return &sonar
}

func (sonar *Sonar) getQaIssueTree(projectKey string, measureType MeasureType) ([]*apistructs.TestIssuesTree, error) {
	sonarIssueTree, err := sonar.invokeSonarIssuesTree(projectKey, measureType)
	if err != nil {
		return nil, err
	}

	var issueTrees []*apistructs.TestIssuesTree

	for _, component := range sonarIssueTree.Components {
		var qaIssueTree apistructs.TestIssuesTree
		qaIssueTree.Path = component.Path
		qaIssueTree.Name = component.Name
		qaIssueTree.Language = component.Language.String()
		qaIssueTree.Measures = component.Measures

		sourceLines, err := sonar.querySonarSourceLinesRendered(component.Key, 1, -1)
		if err != nil {
			return nil, err
		}
		switch measureType {
		case MeasureTypeCoverage:
			for _, line := range sourceLines {
				if line.LineHits != nil {
					qaIssueTree.Lines = append(qaIssueTree.Lines, strconv.Itoa(line.Line))
				}
			}
		case MeasureTypeDuplications:
			for _, line := range sourceLines {
				if line.Duplicated {
					qaIssueTree.Lines = append(qaIssueTree.Lines, strconv.Itoa(line.Line))
				}
			}
		}
		issueTrees = append(issueTrees, &qaIssueTree)
	}

	return issueTrees, nil
}
