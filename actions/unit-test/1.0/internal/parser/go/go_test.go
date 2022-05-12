package _go

import (
	"encoding/json"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/conf"
)

func TestGetGoVersion(t *testing.T) {
	results := getGoVersion()
	t.Log(results)
}

func TestGetSuites(t *testing.T) {
	results, err := getSuites()
	assert.Nil(t, err)
	content, _ := json.Marshal(results)
	t.Log(string(content))
}

func TestGetUtSuites(t *testing.T) {
	suites, err := getUtSuites("golang")
	assert.Nil(t, err)
	content, _ := json.Marshal(suites)
	t.Log(string(content))
}

func TestGoTest(t *testing.T) {
	cfg := conf.Conf{}
	cfg.GoDir = "/Users/ddy/go/src/terminus.io/dice/dice/internal"

	suite, _, err := GoTest("")
	assert.Nil(t, err)

	content, err := json.Marshal(suite)
	assert.Nil(t, err)
	t.Log(string(content))
}
