package build

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/contrast-security/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
)

type ConstrastSecurityResponse struct {
	Count  int           `json:"count"`
	Traces []interface{} `json:"traces"`
}

func Execute() error {
	var cfg conf.Conf
	envconf.MustLoad(&cfg)
	logrus.SetOutput(os.Stdout)
	return build(cfg)
}

func build(cfg conf.Conf) error {
	if len(cfg.Severities) == 0 {
		cfg.Severities = []string{"MEDIUM", "HIGH", "CRITICAL"}
	}
	tr := &http.Transport{
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://app.contrastsecurity.com/Contrast/api/ng/%s/traces/%s/filter", cfg.OrgID, cfg.AppID), nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("expand", cfg.Expand)
	q.Add("status", cfg.Status)
	q.Add("limit", "999")
	for _, severity := range cfg.Severities {
		q.Add("severities", severity)
	}
	auth := fmt.Sprintf("%s:%s", cfg.Username, cfg.ServiceKey)
	auth = base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", auth)
	req.Header.Set("API-Key", cfg.ApiKey)
	req.URL.RawQuery = q.Encode()
	logrus.Info("Contrast Security doing api")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Doing api error: %v", err)
		return err
	}
	var rsp ConstrastSecurityResponse
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Reading body error: %v", err)
		return err
	}
	if err := json.Unmarshal(bytes, &rsp); err != nil {
		logrus.Errorf("Unmarshal response body error: %v", err)
		return err
	}
	if err := storeMetaFile(&cfg, rsp); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully upload metafile\n")
	fmt.Fprintf(os.Stdout, "action meta: count=%d\n", rsp.Count)
	fmt.Fprintf(os.Stdout, "action meta: traces=%s\n", jsonparse.JsonOneLine(rsp.Traces))
	if cfg.AssertCount == 0 {
		return nil
	}
	if rsp.Count > cfg.AssertCount {
		return fmt.Errorf("assert count failed, expected count: %d, actual count: %d", cfg.AssertCount, rsp.Count)
	}
	return nil
}

func storeMetaFile(cfg *conf.Conf, rsp ConstrastSecurityResponse) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "count",
				Value: strconv.FormatInt(int64(rsp.Count), 10),
			},
			{
				Name:  "traces",
				Value: jsonparse.JsonOneLine(rsp.Traces),
			},
		},
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(cfg.MetaFile, string(b), 0644); err != nil {
		return fmt.Errorf("write file:metafile failed")
	}
	return nil
}
