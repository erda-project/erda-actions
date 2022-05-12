package testing

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/integration-test/1.0/internal/conf"
	"github.com/erda-project/erda-proto-go/dop/qa/unittest/pb"
	"github.com/erda-project/erda/pkg/qaparser"
	"github.com/erda-project/erda/pkg/qaparser/surefilexml"
	"github.com/erda-project/erda/pkg/qaparser/testngxml"
)

func MavenTest(cfg *conf.Conf) (*pb.TestSuite, error) {
	var (
		context string
		runCmd  string
		err     error
		mem     float64
		mvnOpts string
		suite   = &pb.TestSuite{
			Totals: &pb.TestTotal{
				Statuses: make(map[string]int64),
			},
			Extra: make(map[string]string),
		}
	)

	// get maven options
	mem = 2048
	if mem, err = strconv.ParseFloat(cfg.PipelineLimitedMem, 64); err != nil {
		logrus.Warning(err)
	}

	if mem > 32 {
		mvnOpts = fmt.Sprintf("-Xmx%gm", mem-32)
	} else {
		mvnOpts = fmt.Sprintf("-Xmx%sm", cfg.PipelineLimitedMem)
	}

	if cfg.RunCmd != "" {
		runCmd = cfg.RunCmd + " || true"
	} else {
		runCmd = fmt.Sprintf("MAVEN_OPTS=%s mvn test -Dmaven.test.failure.ignore=true", mvnOpts)
	}

	context = cfg.Context
	// build
	testCommand := exec.Command("/bin/sh", "-c", runCmd)
	testCommand.Dir = context
	testCommand.Stdout = os.Stderr
	testCommand.Stderr = os.Stderr

	logrus.Infof("Run context=%s, cmd: %s", context, strings.Join(testCommand.Args, " "))

	if err = testCommand.Run(); err != nil {
		return nil, err
	}

	switch strings.ToLower(cfg.ParserType) {
	case TestNg:
		suite, err = getItSuites(TestNgFile, TestNg)
	case Junit:
		suite, err = getItSuites(JunitFile, Junit)
	case "":
		// Get TestNg Suites
		suite, err = getItSuites(TestNgFile, TestNg)
		if suite == nil {
			// Get Junit Suites
			suite, err = getItSuites(JunitFile, Junit)
			if suite == nil {
				return nil, err
			}
		}
	default:
		return nil, errors.Errorf("not support test type:%s", cfg.ParserType)
	}

	if err != nil {
		return nil, err
	}

	if suite == nil {
		return nil, errors.New("nil it results")
	}

	return suite, nil
}

func getItSuites(testFile, testType string) (*pb.TestSuite, error) {
	var (
		suitesParse []*pb.TestSuite
		err         error
	)

	suite := &pb.TestSuite{
		Name: "it-results",
		Totals: &pb.TestTotal{
			Statuses: make(map[string]int64),
		},
		Extra: make(map[string]string),
	}

	files := getFilesPath(".", testFile)
	if len(files) == 0 {
		return nil, errors.Errorf("failed to get it results, file: %s", testFile)
	}

	for _, f := range files {
		if suitesParse, err = getSuites(f, testType); err != nil {
			return nil, errors.Wrapf(err, "failed to parse, file: %s", f)
		}

		for _, s := range suitesParse {
			suite.Tests = append(suite.Tests, s.Tests...)
			totals := &qaparser.Totals{suite.Totals}
			suite.Totals = totals.Add(s.Totals).TestTotal
		}
		suite.Properties = suitesParse[0].Properties
		suite.Package = suitesParse[0].Package
	}

	suite.Extra = setSuiteExtraInfo()

	return suite, nil
}

func setSuiteExtraInfo() map[string]string {
	var (
		// uname -p or uname -s or uname -r
		osArch    = "p"
		osVersion = "r"
		osName    = "s"
	)
	return map[string]string{
		ExtraOsArch:      getOsInfo(osArch),
		ExtraOsVersion:   getOsInfo(osVersion),
		ExtraOsName:      getOsInfo(osName),
		ExtraMvnVersion:  getMvnVersion(),
		ExtraJavaVmName:  getJavaVmVersionInfo(),
		ExtraJavaVersion: getJavaVersionInfo(),
	}
}

func getSuites(f string, testType string) ([]*pb.TestSuite, error) {
	var (
		data   []byte
		err    error
		suites []*pb.TestSuite
		testNg *testngxml.NgTestResult
	)

	if data, err = ioutil.ReadFile(f); err != nil {
		return nil, err
	}

	switch testType {
	case TestNg:
		if testNg, err = testngxml.Ingest(data); err != nil {
			return nil, err
		}

		return testNg.Transfer()
	case Junit:
		if suites, err = surefilexml.Ingest(data); err != nil {
			return nil, err
		}

		return suites, nil
	}

	return nil, errors.Errorf("failed to parse, f: %s, type: %s", f, testType)
}

func getFilesPath(path, prefix string) []string {
	var results []string
	filepath.Walk(path, func(fullPath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(f.Name(), ".xml") &&
			strings.HasPrefix(f.Name(), prefix) &&
			f.Name() != "TEST-TestSuite.xml" {
			results = append(results, fullPath)
		}
		return nil
	})
	return results
}

func getMvnVersion() string {
	var (
		output string
		err    error
	)
	if output, err = runCmd("mvn -v"); err != nil {
		logrus.Errorf("get info by cmd[mvn -v] failed.err=%v", err)
		return ""
	}
	// like this: Apache Maven 3.5.4 (1edded0938998edf8bf061f1ceb3cfdeccf443fe; 2018-06-17T18:33:14Z)
	return strings.SplitN(string(output), " (", 2)[0]
}

func getJavaVersionInfo() string {
	var (
		output string
		err    error
	)
	if output, err = runCmd("java -version"); err != nil {
		logrus.Errorf("get info by cmd[java -version] failed.err=%v", err)
		return ""
	}
	ary := strings.Split(string(output), "\n")
	if ary == nil || len(ary) == 0 {
		return ""
	}

	r := regexp.MustCompile("\"\\d.+\"")
	matches := r.FindString(ary[0])
	if strings.Contains(matches, "\"") {
		matches = strings.Trim(matches, "\"")
	}
	return matches
}

func getJavaVmVersionInfo() string {
	var (
		output string
		err    error
	)
	if output, err = runCmd("java -version"); err != nil {
		logrus.Errorf("get info by cmd[java -version] failed.err=%v", err)
		return ""
	}
	ary := strings.Split(string(output), "\n")
	if ary == nil || len(ary) < 3 {
		logrus.Errorf("get info by cmd[java -version] failed.err=not enough length")
		return ""
	}

	r := ary[2]
	if strings.Contains(r, "(") {
		r = strings.SplitN(r, "(", 2)[0]
	}

	return r
}

func runCmd(cmd string) (string, error) {
	var (
		output []byte
		err    error
	)
	if output, err = exec.Command("/bin/sh", "-c", cmd).CombinedOutput(); err != nil {
		return "", err
	}
	return string(output), nil
}

func getOsInfo(key string) string {
	var (
		output string
		err    error
	)
	cmd := "uname -" + key
	if output, err = runCmd(cmd); err != nil {
		logrus.Errorf("get info by cmd[uname -v] failed.err=%v", err)
		return ""
	}
	return string(output)
}
