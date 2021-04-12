package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/erda-project/erda-actions/actions/mysql-cli/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/httpclient"
)

const (
	mysqlHost     = "MYSQL_HOST"
	mysqlPassword = "MYSQL_PASSWORD"
	mysqlPort     = "MYSQL_PORT"
	mysqlUsername = "MYSQL_USERNAME"
)

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	if err := build(cfg); err != nil {
		return err
	}

	return nil
}

func build(cfg conf.Conf) error {

	err := simpleRunAndPrint("/bin/sh", "-c", "echo '"+cfg.Sql+"' >> /mysql-cli.sql")
	if err != nil {
		return err
	}
	mysqlAddon, err := getAddonFetchResponseData(cfg)
	if err != nil {
		fmt.Println(fmt.Errorf("getAddonFetchResponseData error %v", err))
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

	cmd := exec.Command("mysql", "-h"+mysqlAddon.Config[mysqlHost].(string), "-P"+mysqlAddon.Config[mysqlPort].(string),
		"-u"+mysqlAddon.Config[mysqlUsername].(string), "-p"+mysqlAddon.Config[mysqlPassword].(string), "-D"+cfg.Database,
		"-e", "source /mysql-cli.sql")

	var output bytes.Buffer
	var errors bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &errors
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Errorf("exec sql error %v", err))
	}

	// error 信息大于 0
	if errors.Len() > 0 {
		return fmt.Errorf("exec sql error: %s", errors.String())
	}

	split := strings.Split(output.String(), "\n")
	for index, v := range split {
		err = simpleRun("/bin/sh", "-c", "echo 'exec_result["+strconv.Itoa(index)+"]="+v+"'>> "+cfg.MetaFile)
		if err != nil {
			return err
		}
	}

	fmt.Println(output.String())
	return nil
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
