package js

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/base"
	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/parser"
	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/tap"
	"github.com/erda-project/erda-proto-go/dop/qa/unittest/pb"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/qaparser"
)

func JsTest(codePath string) (*pb.TestSuite, []*pb.CodeCoverageNode, error) {
	var (
		content []byte
		err     error
	)
	if err = base.ChangeWorkDir(codePath); err != nil {
		return nil, nil, err
	}

	// install npm mocha
	if err := base.ExecuteCmd("npm install mocha-tap-reporter; npm install"); err != nil {
		return nil, nil, err
	}

	if content, err = base.ExecuteCmdOutput("npm run test | grep '^[^>]'"); err != nil {
		return nil, nil, err
	}
	suites, err := getUtSuites(string(content), codePath)
	if err != nil {
		return nil, nil, err
	}

	return suites, nil, nil
}

func getUtSuites(content, moduleName string) (*pb.TestSuite, error) {
	var (
		suiteTap *tap.Testsuite
		err      error
	)
	if suiteTap, err = parser.ParserTapSuite(content); err != nil {
		return nil, err
	}

	suite := &pb.TestSuite{
		Name: fmt.Sprintf("js(ut):%s", moduleName),
		Totals: &pb.TestTotal{
			Statuses: make(map[string]int64),
		},
		Extra: make(map[string]string),
	}

	for _, t := range suiteTap.Tests {
		tests := &pb.Test{}
		if !t.Ok {
			tests.Status = "failed"
			tests.Error = &pb.TestError{
				Body:    t.Diagnostic,
				Message: "Failed",
			}
		} else {
			tests.Status = "passed"
			if t.Diagnostic != "" {
				makeTotals(t.Diagnostic, suite.Totals)
			}
		}
		tests.Stdout = t.Description
		tests.Name = t.Description
		suite.Tests = append(suite.Tests, tests)
	}
	suite.Extra = setSuiteExtraInfo()

	return suite, nil
}

func makeTotals(diagnostic string, totals *pb.TestTotal) {
	var (
		val int
		err error
	)

	if totals.Statuses == nil {
		totals.Statuses = qaparser.NewStatuses(0, 0, 0, 0)
	}

	if diagnostic != "" {
		sl := strings.Split(diagnostic, "\n")
		for _, s := range sl {
			ret := strings.Split(s, " ")
			if len(ret) != 2 {
				logrus.Warnf("split error: %+v", s)
				continue
			}

			if val, err = strconv.Atoi(ret[1]); err != nil {
				logrus.Errorf("failed to convert, value: %d, (%+v)", ret[1], err)
			}
			dat := int64(val)

			switch ret[0] {
			case "tests":
				totals.Tests += dat
			case "pass":
				totals.Statuses[string(apistructs.TestStatusPassed)] += dat
			case "fail":
				totals.Statuses[string(apistructs.TestStatusFailed)] += dat
			case "skip":
				totals.Statuses[string(apistructs.TestStatusSkipped)] += dat
			}
		}
	}
}

func setSuiteExtraInfo() map[string]string {
	var (
		// uname -p or uname -s or uname -r
		osArch    = "p"
		osVersion = "r"
		osName    = "s"
	)
	return map[string]string{
		base.ExtraOsArch:    base.GetOsInfo(osArch),
		base.ExtraOsVersion: base.GetOsInfo(osVersion),
		base.ExtraOsName:    base.GetOsInfo(osName),
		"node.version":      getJsVersion(),
	}
}

func getJsVersion() string {
	var (
		output string
		err    error
	)
	if output, err = base.RunCmd("node --version"); err != nil {
		logrus.Warningf("failed to get node version, (%+v)", err)
		return ""
	}

	return output
}
