package testing

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/erda-project/erda-actions/actions/integration-test/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

func TestExec(t *testing.T) {
	var cfg conf.Conf
	err := envconf.Load(&cfg)
	require.NoError(t, err)
	cfg.RunCmd = "/usr/local/Cellar/maven/3.5.2/libexec/bin/mvn test -Dmaven.test.failure.ignore=true"
	Exec(&cfg)
}

func TestFileGlobal(t *testing.T) {
	list, err := filepath.Glob("TEST-*.xml")
	assert.Nil(t, err)
	logrus.Info(list)
}

func TestSplit(t *testing.T) {
	var tmp = filepath.Join(filepath.Join("root", "target/surefire-reports", "/TEST-1.xml"))
	logrus.Info(tmp[strings.LastIndex(tmp, "/")+1:])
	logrus.Info(fmt.Sprintf("%d-%s", time.Now().Unix(), strings.Replace(tmp, "/", "-", -1)))
}

func TestCmd(t *testing.T) {
	output, err := exec.Command("/bin/bash", "-c", "java -version").CombinedOutput()
	assert.Nil(t, err)
	logrus.Info(string(output))

	str := string(output)
	logrus.Info(str)

	ary := strings.Split(str, "\n")
	regex := "\"\\d.+\""
	r := regexp.MustCompile(regex)

	matches := r.FindString(ary[0])

	logrus.Info(strings.Trim(matches, "\""))
	logrus.Info(getJavaVmVersionInfo())
	logrus.Info(getJavaVersionInfo())
}
