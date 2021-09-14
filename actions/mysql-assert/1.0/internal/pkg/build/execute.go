package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/mysql-assert/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/assert"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

const (
	mysqlHost     = "MYSQL_HOST"
	mysqlPassword = "MYSQL_PASSWORD"
	mysqlPort     = "MYSQL_PORT"
	mysqlUsername = "MYSQL_USERNAME"
)

func Execute() error {
	var cfg conf.Conf
	logrus.SetOutput(os.Stdout)
	envconf.MustLoad(&cfg)

	if err := build(cfg); err != nil {
		return err
	}

	return nil
}

type results struct {
	Rows []interface{} `json:"rows"`
}

func build(cfg conf.Conf) error {

	mysqlAddon, err := getAddonFetchResponseData(cfg)
	if err != nil {
		logrus.Errorf("getAddonFetchResponseData error %v", err)
		return err
	}
	if mysqlAddon == nil {
		return fmt.Errorf("not find this %s mysql service", cfg.DataSource)
	}

	if mysqlAddon.Config[mysqlHost] == nil {
		return fmt.Errorf("not find %s", mysqlHost)
	}
	if mysqlAddon.Config[mysqlPort] == nil {
		return fmt.Errorf("not find %s", mysqlPort)
	}
	if mysqlAddon.Config[mysqlUsername] == nil {
		return fmt.Errorf("not find %s", mysqlUsername)
	}
	if mysqlAddon.Config[mysqlPassword] == nil {
		return fmt.Errorf("not find %s", mysqlPassword)
	}

	logrus.Infof("----------- execute sql ----------- ")
	logrus.Infof(cfg.Sql)

	mysqlFile, err := os.Create("mysql-cli.sql")
	if err != nil {
		return err
	}
	_, err = mysqlFile.Write([]byte(cfg.Sql))
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", "mysqlsh --json=raw --host="+mysqlAddon.Config[mysqlHost].(string)+" --password="+mysqlAddon.Config[mysqlPassword].(string)+" "+
		"--dbuser="+mysqlAddon.Config[mysqlUsername].(string)+" --port="+mysqlAddon.Config[mysqlPort].(string)+" --database="+cfg.Database+" --file=mysql-cli.sql")

	var output bytes.Buffer
	var errors bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &errors
	err = cmd.Run()
	if err != nil {
		logrus.Errorf("exec sql error %v", err)
	}

	// error 信息大于 0
	if errors.Len() > 0 {
		split := strings.Split(errors.String(), "\n")
		for _, v := range split {
			if strings.Contains(v, "[Warning] Using a password on the command line interface can be insecure.") {
				continue
			}
			if strings.Contains(v, "WARNING: The --dbuser option has been deprecated,") {
				continue
			}
			if strings.TrimSpace(v) == "" {
				continue
			}
			logrus.Infof("------ sql exec failed -------")
			return fmt.Errorf("%v", v)
		}
	}

	logrus.Infof("------ sql exec done -------")
	result, err := getResultMeta(output)
	if err != nil {
		return err
	}

	logrus.Infof("------ json parsing -------")
	allOutput, err := parseJson(cfg, result)
	if err != nil {
		return err
	}
	logrus.Infof("------ json parsing done -------")

	logrus.Infof("------ output assert -------")
	allSuccess := assertValue(cfg, allOutput)
	if !allSuccess {
		return fmt.Errorf("asssert faild")
	}
	logrus.Infof("------ output assert done -------")

	allOutput["result"] = result
	err = storeMetaFile(&cfg, allOutput)
	if err != nil {
		return err
	}
	return nil
}

func assertValue(cfg conf.Conf, actualValues map[string]string) bool {
	var allSuccess = true
	for _, v := range cfg.OutParams {

		if len(v.Assert) <= 0 {
			continue
		}

		success, err := assert.DoAssert(actualValues[v.Key], v.Assert, jsonparse.JsonOneLine(v.Value))
		if err != nil || !success {
			allSuccess = false
		}
		// to assert
		logrus.Infof("Assert Result:")
		logrus.Infof("  value: %v", jsonparse.JsonOneLine(v.Value))
		logrus.Infof("  assert: %v", v.Assert)
		logrus.Infof("  actualValue: %s", actualValues[v.Key])
		logrus.Infof("  success: %v", success)
		logrus.Infof("==========")
	}

	logrus.Infof("AllAssert Result: %v", allSuccess)
	return allSuccess
}

func parseJson(cfg conf.Conf, data string) (map[string]string, error) {

	var out bytes.Buffer
	err := json.Indent(&out, []byte(data), "", "\t")
	logrus.Infof("json data:")
	if err != nil {
		logrus.Infof("%s\n", out.String())
	} else {
		logrus.Infof("%s\n", data)
	}
	var allOutput = make(map[string]string)

	for _, express := range cfg.OutParams {

		if len(express.Key) <= 0 || len(express.Expression) <= 0 {
			continue
		}

		result := jsonparse.FilterJson([]byte(data), express.Expression, apistructs.APIOutParamSourceBodyJson.String())
		logrus.Infof("Out Params:")
		logrus.Infof("  key: %v", express.Key)
		logrus.Infof("  expr: %v", express.Expression)
		logrus.Infof("  value: %v", jsonparse.JsonOneLine(result))
		logrus.Infof("==========")
		allOutput[express.Key] = jsonparse.JsonOneLine(result)
	}

	return allOutput, nil
}

func getResultMeta(output bytes.Buffer) (string, error) {
	split := strings.Split(strings.TrimSpace(output.String()), "\n")
	v := split[len(split)-1]
	var result results
	err := json.Unmarshal([]byte(v), &result)
	if err != nil {
		logrus.Errorf("unmarshal result error: %v", err)
		printJson(v)
		return "", err
	} else {
		rows, err := json.Marshal(result.Rows)
		if err != nil {
			logrus.Errorf("marshal rows error: %v", err)
			printJson(v)
			return "", err
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, rows, "", "\t")
		if err != nil {
			logrus.Errorf("format rows result error: %v", err)
			logrus.Infof(string(rows))
			return "", err
		}
		logrus.Infof(prettyJSON.String())

		return prettyJSON.String(), nil
	}
}

func printJson(v string) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(v), "", "\t")
	if err != nil {
		logrus.Errorf("format result error: %v", err)
		logrus.Infof(v)
	} else {
		logrus.Infof("sql select json: ", prettyJSON.String())
	}
}

func simpleRun(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func simpleRunAndPrint(name string, arg ...string) error {
	fmt.Fprintf(os.Stdout, "Run: %s, %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func storeMetaFile(cfg *conf.Conf, results map[string]string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{},
	}

	for key, value := range results {
		{
			meta.Metadata = append(meta.Metadata, apistructs.MetadataField{
				Name:  key,
				Value: value,
			})
		}
	}

	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(cfg.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func getAddonFetchResponseData(cfg conf.Conf) (*apistructs.AddonFetchResponseData, error) {
	var buffer bytes.Buffer
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(cfg.DiceOpenapiAddr).
		Header("Authorization", cfg.DiceOpenapiToken).
		Path(fmt.Sprintf("/api/addons/%s", cfg.DataSource)).
		Do().Do().Body(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to getAddonFetchResponseData, err: %v", err)
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("failed to getAddonFetchResponseData, statusCode: %d, respBody: %s", resp.StatusCode(), buffer.String())
	}
	var result apistructs.AddonFetchResponse
	respBody := buffer.String()
	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return nil, fmt.Errorf("failed to getAddonFetchResponseData, err: %v, json string: %s", err, respBody)
	}
	return &result.Data, nil
}
