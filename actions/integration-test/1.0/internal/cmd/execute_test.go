package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/erda-project/erda-actions/actions/integration-test/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

func TestReplaceM2File(t *testing.T) {
	os.Setenv("BP_NEXUS_URL", "http://nexus.url")
	os.Setenv("BP_NEXUS_USERNAME", "username")
	os.Setenv("BP_NEXUS_PASSWORD", "password")

	var cfg conf.Conf
	err := envconf.Load(&cfg)
	require.NoError(t, err)

	err = replaceM2File(&cfg)
	assert.Nil(t, err)
}
