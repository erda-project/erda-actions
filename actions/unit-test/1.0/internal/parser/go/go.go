package _go

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/base"
	"github.com/erda-project/erda-proto-go/dop/qa/unittest/pb"
	"github.com/erda-project/erda/pkg/qaparser"
	"github.com/erda-project/erda/pkg/qaparser/surefilexml"
)

func GoTest(codePath string) (*pb.TestSuite, []*pb.CodeCoverageNode, error) {
	resultFile := fmt.Sprintf("%s-%s", base.Golang, base.Cfg.GittarCommit)

	if base.Cfg.GoDir == "" {
		return nil, nil, errors.New("need params go_dir")
	}

	var (
		contextDir string
		err        error
	)
	contextDir = codePath

	goWorkSpace := filepath.Join("/opt", base.Cfg.GoDir)
	// make go dir and copy code to goPath.
	err = base.ExecuteCmd(fmt.Sprintf("mkdir -p %s; cp -rf %s %s",
		filepath.Dir(goWorkSpace), contextDir, goWorkSpace))
	if err != nil {
		return nil, nil, err
	}

	if base.Cfg.Command != "" {
		if err := base.ExecuteCmd(base.Cfg.Command); err != nil {
			return nil, nil, err
		}
	}

	if err := base.ExecuteCmd("go mod tidy"); err != nil {
		return nil, nil, err
	}
	testCmd := fmt.Sprintf("go test -v ./... | go-junit-report > %s", resultFile)
	if err := base.ExecuteCmd(testCmd); err != nil {
		return nil, nil, err
	}

	suites, err := getUtSuites(codePath)
	if err != nil {
		return nil, nil, err
	}

	return suites, nil, nil
}

func getUtSuites(moduleName string) (*pb.TestSuite, error) {
	var (
		suitesParse []*pb.TestSuite
		err         error
	)

	suite := &pb.TestSuite{
		Name: fmt.Sprintf("golang(ut):%s", moduleName),
		Totals: &pb.TestTotal{
			Statuses: make(map[string]int64),
		},
		Extra: make(map[string]string),
	}

	if suitesParse, err = getSuites(); err != nil {
		return nil, errors.Wrapf(err, "failed to Parse golang xml file to suites")
	}
	if len(suitesParse) == 0 {
		return nil, errors.New("suites is nil")
	}

	for _, s := range suitesParse {
		suite.Tests = append(suite.Tests, s.Tests...)
		totals := &qaparser.Totals{suite.Totals}
		suite.Totals = totals.Add(s.Totals).TestTotal
	}
	suite.Properties = suitesParse[0].Properties
	suite.Package = suitesParse[0].Package

	suite.Extra = setSuiteExtraInfo()

	return suite, nil
}

func getSuites() ([]*pb.TestSuite, error) {
	var (
		data   []byte
		err    error
		suites []*pb.TestSuite
	)

	resultFile := fmt.Sprintf("%s-%s", base.Golang, base.Cfg.GittarCommit)
	_, err = os.Stat(resultFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(err, fmt.Sprintf("%s not exist", resultFile))
		}
		return nil, errors.Wrap(err, fmt.Sprintf("%s does exist, but throw other error when checking", resultFile))
	}
	if data, err = ioutil.ReadFile(resultFile); err != nil {
		return nil, err
	}

	if suites, err = surefilexml.Ingest(data); err != nil {
		return nil, err
	}

	return suites, nil
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
		ExtraGoVersion:      getGoVersion(),
	}
}

func getGoVersion() string {
	var (
		output string
		err    error
	)

	if output, err = base.RunCmd("go version"); err != nil {
		logrus.Warningf("failed to get go version, (%+v)", err)
		return ""
	}

	ary := strings.Split(string(output), " ")
	if ary == nil || len(ary) < 3 {
		return ""
	}

	return ary[2]
}
