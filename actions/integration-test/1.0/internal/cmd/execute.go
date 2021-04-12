package cmd

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/integration-test/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/integration-test/1.0/internal/testing"
	"github.com/erda-project/erda/pkg/envconf"
)

/*
	Input:
		1. stdin: {"source": {...}, "params":{...}}
		2. args[1]: directory

	Output:
		stdout: {"version":{...},metadata:{...}}
*/

var (
	configMap = map[string]string{}
	m2File    = "/root/.m2/settings.xml"
)

func Execute() error {
	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}

	// replace m2 file
	if err := replaceM2File(&cfg); err != nil {
		return err
	}

	if err := testing.Exec(&cfg); err != nil {
		return errors.Wrapf(err, "execute testing failed. error=[%v]", err)
	}

	return nil
}

func replaceM2File(cfg *conf.Conf) error {
	configMap["BP_NEXUS_URL"] = "http://" + strings.TrimPrefix(cfg.NexusUrl, "http://")
	configMap["BP_NEXUS_USERNAME"] = cfg.NexusUsername
	configMap["BP_NEXUS_PASSWORD"] = cfg.NexusPassword

	if bytes, err := ioutil.ReadFile(m2File); err == nil {
		result, _ := renderConfig(string(bytes))
		ioutil.WriteFile(m2File, []byte(result), os.ModePerm)
	} else {
		return errors.Errorf("read maven file fail %v", err)
	}
	return nil
}

func renderConfig(template string) (string, bool) {
	compile, _ := regexp.Compile("{{.+?}}")
	hasChange := false
	result := compile.ReplaceAllStringFunc(template, func(s string) string {
		key := s[2:(len(s) - 2)]
		value, ok := configMap[key]
		if ok {
			hasChange = true
			return value
		} else {
			return s
		}
	})
	return result, hasChange
}
