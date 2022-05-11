package ut

import (
	"encoding/json"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	_go "github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/parser/go"
	"github.com/erda-project/erda-proto-go/dop/qa/unittest/pb"
)

func TestCheckLanguage(t *testing.T) {
	lan, err := checkLanguage(".")
	assert.Nil(t, err)
	t.Log(lan)
}

func TestGetResults(t *testing.T) {
	var suites []*pb.TestSuite

	suite, _, err := _go.GoTest("")
	assert.Nil(t, err)

	suites = append(suites, suite)

	result, err := makeUtResults(suites)
	assert.Nil(t, err)

	c, err := json.Marshal(result)
	assert.Nil(t, err)

	t.Log(string(c))
}
