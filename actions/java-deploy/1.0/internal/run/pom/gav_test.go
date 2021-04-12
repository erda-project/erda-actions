package pom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGAV(t *testing.T) {
	gav, err := GetGAV("./testdata/pom.xml")
	assert.NoError(t, err)
	assert.NotNil(t, gav)
	assert.Equal(t, "io.terminus", gav.GroupID)
	assert.Equal(t, "dice-test", gav.ArtifactID)
	assert.Equal(t, "1.0-SNAPSHOT", gav.Version)
}
