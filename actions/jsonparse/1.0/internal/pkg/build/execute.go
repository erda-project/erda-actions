package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/erda-project/erda-actions/actions/jsonparse/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/envconf"
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
	var out bytes.Buffer
	err := json.Indent(&out, []byte(cfg.Data), "", "\t")
	fmt.Printf("json data:")
	if err != nil {
		fmt.Printf("%s\n", out.String())
	} else {
		fmt.Printf("%s\n", cfg.Data)
	}
	for _, express := range cfg.OutParams {
		result := jsonparse.FilterJson([]byte(cfg.Data), express.Expression, apistructs.APIOutParamSourceBodyJson.String())
		fmt.Printf("Out Params:")
		fmt.Printf("  key: %v", express.Key)
		fmt.Printf("  expr: %v", express.Expression)
		fmt.Printf("  value: %v", jsonparse.JsonOneLine(result))
		fmt.Printf("==========")
		err := simpleRun("/bin/sh", "-c", "echo '"+express.Key+"="+jsonparse.JsonOneLine(result)+"'>> "+cfg.MetaFile)
		if err != nil {
			return fmt.Errorf("echod result error: %v", err)
		}
	}
	return nil
}

func simpleRun(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
