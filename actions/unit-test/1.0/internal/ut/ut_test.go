package ut

import (
	"encoding/json"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/conf"
	_go "github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/parser/go"
	"github.com/erda-project/erda/apistructs"
)

func TestCheckLanguage(t *testing.T) {
	lan, err := checkLanguage(".")
	assert.Nil(t, err)
	t.Log(lan)
}

func TestGetResults(t *testing.T) {
	cfg := conf.Conf{}
	cfg.GoDir = "/Users/ddy/go/src/terminus.io/dice/dice/internal"

	var suites []*apistructs.TestSuite

	suite, err := _go.GoTest("")
	assert.Nil(t, err)

	suites = append(suites, suite)

	result, err := makeUtResults(suites)
	assert.Nil(t, err)

	c, err := json.Marshal(result)
	assert.Nil(t, err)

	t.Log(string(c))
}
