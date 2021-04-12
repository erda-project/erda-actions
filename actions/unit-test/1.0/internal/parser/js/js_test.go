package js

import (
	"encoding/json"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

var tapResult = `1..3
ok 1 A suite contains spec with an expectation
ok 2 A suite contains spec with an expectation
ok 3 A suite contains spec with an expectation
# tests 3
# pass 3
# fail 0
# skip 0`

func TestGetJsVersion(t *testing.T) {
	results := getJsVersion()
	t.Log(results)
}

func TestGetUtSuites(t *testing.T) {
	suites, err := getUtSuites(tapResult, "test")
	assert.Nil(t, err)

	c, err := json.Marshal(suites)
	assert.Nil(t, err)
	t.Log(string(c))
}

func TestJsTest(t *testing.T) {
	suite, err := JsTest("")
	assert.Nil(t, err)

	content, err := json.Marshal(suite)
	assert.Nil(t, err)
	t.Log(string(content))
}
