package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/erda-project/erda-actions/actions/redis-cli/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

const (
	redisHost     = "REDIS_HOST"
	redisPassword = "REDIS_PASSWORD"
	redisPort     = "REDIS_PORT"
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
	mysqlFile, err := os.Create("/redis-cli.txt")
	if err != nil {
		return err
	}
	_, err = mysqlFile.Write([]byte(cfg.Command))
	if err != nil {
		return err
	}

	redisAddon, err := getAddonFetchResponseData(cfg)
	if err != nil {
		return fmt.Errorf("getAddonFetchResponseData error %v", err)
	}
	if redisAddon == nil {
		return fmt.Errorf("not find this %s mysql service", cfg.DataSource)
	}

	if redisAddon.Config[redisHost] == nil {
		return fmt.Errorf("not find %s", redisHost)
	}
	if redisAddon.Config[redisPort] == nil {
		return fmt.Errorf("not find %s", redisPort)
	}
	if redisAddon.Config[redisPassword] == nil {
		return fmt.Errorf("not find %s", redisPassword)
	}

	fmt.Fprintf(os.Stdout, "Run: %s, %s\n", "/bin/sh -c", fmt.Sprint("cat /redis-cli.txt | redis-cli -h "+redisAddon.Config[redisHost].(string)+" -p "+redisAddon.Config[redisPort].(string)+" -a "+redisAddon.Config[redisPassword].(string)))
	cmd := exec.Command("/bin/sh", "-c", "cat /redis-cli.txt | redis-cli -h "+redisAddon.Config[redisHost].(string)+" -p "+redisAddon.Config[redisPort].(string)+" -a "+redisAddon.Config[redisPassword].(string))

	var output []byte
	if output, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cmd exec output error: %v print: %s", err, string(output))
	}

	split := strings.Split(string(output), "\n")
	for index, v := range split {
		err = simpleRun("/bin/sh", "-c", "echo 'exec_result["+strconv.Itoa(index)+"]="+v+"'>> "+cfg.MetaFile)
		if err != nil {
			return fmt.Errorf("print result error: %v", err)
		}
	}

	fmt.Println(string(output))
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
		Path(fmt.Sprintf("/api/addons/%s", cfg.DataSource)).
		Header("Authorization", cfg.DiceOpenapiToken).
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
