package main

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	"github.com/erda-project/erda/pkg/crypto/uuid"
)

func TestSonarCreateProject(t *testing.T) {
	projectKey := uuid.UUID()
	err := sonar.createProject(projectKey)
	assert.NoError(t, err)
}
