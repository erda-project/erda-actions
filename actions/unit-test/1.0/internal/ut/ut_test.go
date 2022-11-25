package ut

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	_go "github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/parser/go"
	"github.com/erda-project/erda-proto-go/dop/qa/unittest/pb"
	"github.com/erda-project/erda/apistructs"
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

func Test_calculateTestCoverage(t *testing.T) {
	type args struct {
		coverage *pb.TestCallBackRequest
	}
	tests := []struct {
		name        string
		args        args
		passedRate  string
		failedRate  string
		skippedRate string
	}{
		{
			name: "test calculateTestCoverage",
			args: args{
				coverage: &pb.TestCallBackRequest{
					Totals: &pb.TestTotal{
						Statuses: map[string]int64{
							string(apistructs.TestStatusPassed):  2,
							string(apistructs.TestStatusFailed):  2,
							string(apistructs.TestStatusSkipped): 2,
							string(apistructs.TestStatusError):   2,
						},
					},
				},
			},
			passedRate:  "25.00",
			failedRate:  "50.00",
			skippedRate: "25.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passedRate, failedRate, skippedRate := calculateTestCoverage(tt.args.coverage)
			assert.Equal(t, tt.passedRate, fmt.Sprintf("%.2f", passedRate))
			assert.Equal(t, tt.failedRate, fmt.Sprintf("%.2f", failedRate))
			assert.Equal(t, tt.skippedRate, fmt.Sprintf("%.2f", skippedRate))
		})
	}
}
