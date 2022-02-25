package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSonarQuerySonarProjectMeasures(t *testing.T) {
	measures, err := sonar.querySonarProjectMeasures("57ae043d89274eec9431b78c2ab71954")
	assert.NoError(t, err)
	statistics := makeQATestIssueStatistics(*measures)
	fmt.Println(statistics.Rating)
}
