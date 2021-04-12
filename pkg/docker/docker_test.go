package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	err := Login("docker-hosted-nexus-sys.dev.terminus.io", "admin", "admin123")
	assert.NoError(t, err)
}
