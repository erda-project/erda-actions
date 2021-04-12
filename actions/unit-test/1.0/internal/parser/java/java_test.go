package java

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestGetFilePath(t *testing.T) {
	results := getFilesPath(".", "testng-results.xml")
	t.Log(results)
}

func TestGetMvnVersion(t *testing.T) {
	results := getMvnVersion()
	t.Log(results)
}

func TestGetJavaVersionInfo(t *testing.T) {
	results := getJavaVersionInfo()
	t.Log(results)
}

func TestGetJavaVmVersionInfo(t *testing.T) {
	results := getJavaVmVersionInfo()
	t.Log(results)
}

func TestGetSuites(t *testing.T) {
	// testng
	results, err := getSuites("testdata/target/surefire-reports/testng-results.xml", TestNg)
	assert.Nil(t, err)
	fmt.Printf("testng results:%+v\n", results)

	// junit
	results, err = getSuites("testdata/target/surefire-reports/TEST-io.terminus.dice.project.appci.CicdTest.xml", Junit)
	assert.Nil(t, err)
	fmt.Printf("junit results:%+v\n", results)
}

func TestGetUtSuites(t *testing.T) {
	suites := getUtSuites(TestNgFile, TestNg, "testng")
	content, err := json.Marshal(suites)
	assert.Nil(t, err)
	t.Log(string(content))

	suites = getUtSuites(JunitFile, Junit, "junit")
	content, err = json.Marshal(suites)
	assert.Nil(t, err)
	t.Log(string(content))
}

func TestMavenTest(t *testing.T) {
	t.Log(os.Getenv("JAVA_HOME"))
	suite, err := MavenTest("")
	assert.Nil(t, err)

	content, err := json.Marshal(suite)
	assert.Nil(t, err)
	t.Log(string(content))
}
