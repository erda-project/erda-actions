package main

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

var sonar *Sonar

func init() {
	sonar = &Sonar{Auth: SonarAuth{
		HostURL:  "https://sonar-sys.test.terminus.io",
		Login:    "admin",
		Password: "suqing",
	}}
}

func TestSonarCreateQualityGate(t *testing.T) {
	var (
		err  error
		name string
	)
	name = "my-cq-1"
	err = sonar.createQualityGate(name)
	assert.NoError(t, err)

	// too-long-name
	name = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789a"
	err = sonar.createQualityGate(name)
	assert.NoError(t, err)
}

func TestSonarCreateQualityGateCondition(t *testing.T) {
	var (
		err      error
		gateName string
	)

	gateName = "my-cq-1"
	cond := QualityGateCondition{
		Error:  1,
		Metric: MetricKeyNewSecurityRating,
		Op:     OpLT,
	}
	err = sonar.createQualityGateCondition(gateName, cond)
	assert.NoError(t, err)
}
