package java

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	. "github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/base"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/qaparser"
	"github.com/erda-project/erda/pkg/qaparser/surefilexml"
	"github.com/erda-project/erda/pkg/qaparser/testngxml"
)

func MavenTest(codePath string) (*apistructs.TestSuite, error) {
	if err := ChangeWorkDir(codePath); err != nil {
		return nil, err
	}

	// get maven options
	var (
		mem     float64
		err     error
		mvnOpts string
	)

	mem = 2048
	if mem, err = strconv.ParseFloat(Cfg.PipelineLimitedMem, 64); err != nil {
		logrus.Warning(err)
	}

	if mem > 32 {
		mvnOpts = fmt.Sprintf("-Xmx%gm", mem-32)
	} else {
		mvnOpts = fmt.Sprintf("-Xmx%sm", Cfg.PipelineLimitedMem)
	}

	var execCmd string
	if _, err := os.Stat("./pom.xml"); err != nil {
		if os.IsNotExist(err) {
			execCmd = fmt.Sprintf("GRADLE_OPTS=%s ./gradlew test", mvnOpts)
		} else {
			// real error
			return nil, err
		}
	} else {
		// found
		execCmd = fmt.Sprintf("MAVEN_OPTS=%s mvn test -Dmaven.test.failure.ignore=true", mvnOpts)
	}
	if Cfg.Command != "" {
		execCmd = fmt.Sprintf("GRADLE_OPTS=%s MAVEN_OPTS=%s %s", mvnOpts, mvnOpts, Cfg.Command)
	}

	var execError error
	// Ut判断结果文件是否已经存在，UT分析去重
	files := getFilesPath(".", TestNgFile)
	if len(files) == 0 {
		files = getFilesPath(".", JunitFile)
		if len(files) == 0 {
			if execError = ExecuteCmd(execCmd); err != nil {
				logrus.Errorf("exec cmd: %s failed, err: %v", execCmd, err)
			}
		} else {
			logrus.Infof("no need execute mvn test for junit file already exist")
		}
	} else {
		logrus.Infof("no need execute mvn test for testNg file already exist")
	}

	// Get TestNg Suites
	suite := getUtSuites(TestNgFile, TestNg, codePath)
	if suite == nil {
		// Get Junit Suites
		suite = getUtSuites(JunitFile, Junit, codePath)
		if suite == nil {
			return nil, errors.New("nil Suites")
		}
	}

	return suite, execError
}

func getUtSuites(testFile, testType, moduleName string) *apistructs.TestSuite {
	suite := &apistructs.TestSuite{
		Name: fmt.Sprintf("java(ut):%s", moduleName),
		Totals: &apistructs.TestTotals{
			Statuses: make(map[apistructs.TestStatus]int),
		},
		Extra: make(map[string]string),
	}

	files := getFilesPath(".", testFile)
	if len(files) == 0 {
		logrus.Warningf("not exist, file: %s", testFile)
		return nil
	}

	for _, f := range files {
		suitesParse, err := getSuites(f, testType)
		if err != nil {
			logrus.Warningf("failed to parse, file :%s, (%+v)", f, err)
			continue
		}
		if len(suitesParse) == 0 {
			logrus.Warningf("nil suites, file: %s", f)
			continue
		}

		for _, s := range suitesParse {
			suite.Tests = append(suite.Tests, s.Tests...)
			totals := &qaparser.Totals{suite.Totals}
			suite.Totals = totals.Add(s.Totals).TestTotals
		}
		suite.Properties = suitesParse[0].Properties
		suite.Package = suitesParse[0].Package
	}

	suite.Extra = setSuiteExtraInfo()

	return suite
}

func setSuiteExtraInfo() map[string]string {
	var (
		// uname -p or uname -s or uname -r
		osArch    = "p"
		osVersion = "r"
		osName    = "s"
	)
	return map[string]string{
		ExtraOsArch:      GetOsInfo(osArch),
		ExtraOsVersion:   GetOsInfo(osVersion),
		ExtraOsName:      GetOsInfo(osName),
		ExtraMvnVersion:  getMvnVersion(),
		ExtraJavaVmName:  getJavaVmVersionInfo(),
		ExtraJavaVersion: getJavaVersionInfo(),
	}
}

func getSuites(f string, testType string) ([]*apistructs.TestSuite, error) {
	var (
		data   []byte
		err    error
		suites []*apistructs.TestSuite
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

	return nil, nil
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

	if output, err = RunCmd("mvn -v"); err != nil {
		logrus.Warningf("failed to get maven version, (%+v)", err)
		return ""
	}

	// like this: Apache Maven 3.5.4 (1edded0938998edf8bf061f1ceb3cfdeccf443fe; 2018-06-17T18:33:14Z)
	return strings.SplitN(output, " (", 2)[0]
}

func getJavaVersionInfo() string {
	var (
		output string
		err    error
	)

	if output, err = RunCmd("java -version"); err != nil {
		logrus.Warningf("failed to get java version, (%+v)", err)
		return ""
	}

	ary := strings.Split(output, "\n")
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
	if output, err = RunCmd("java -version"); err != nil {
		logrus.Warningf("failed to get java version, (%+v)", err)
		return ""
	}

	ary := strings.Split(output, "\n")
	if ary == nil || len(ary) < 3 {
		logrus.Warning("failed to get java version with not enough length")
		return ""
	}

	r := ary[2]
	if strings.Contains(r, "(") {
		r = strings.SplitN(r, "(", 2)[0]
	}

	return r
}
