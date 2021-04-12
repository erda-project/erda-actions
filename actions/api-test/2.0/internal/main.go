package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"

	"github.com/erda-project/erda-actions/actions/api-test/2.0/internal/cookiejar"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/apitestsv2"
	"github.com/erda-project/erda/pkg/envconf"
)

const CookieJar = "cookieJar"

func main() {
	log.Init()

	// print logo
	printLogo()

	// parse conf from env
	var cfg EnvConfig
	if err := envconf.Load(&cfg); err != nil {
		logrus.Fatalf("failed to parse config from env, err: %v\n", err)
	}

	// get api info and pretty print
	apiInfo := generateAPIInfoFromEnv(cfg)

	printOriginalAPIInfo(apiInfo)

	// success
	var success = true
	defer func() {
		if !success {
			os.Exit(1)
		}
	}()

	// defer create metafile
	meta := NewMeta()
	meta.OutParamsDefine = apiInfo.OutParams
	defer writeMetaFile(cfg.MetaFile, meta)

	// global config
	var apiTestEnvData *apistructs.APITestEnvData
	caseParams := make(map[string]*apistructs.CaseParams)
	if cfg.GlobalConfig != nil {
		apiTestEnvData = &apistructs.APITestEnvData{}
		apiTestEnvData.Domain = cfg.GlobalConfig.Domain
		apiTestEnvData.Header = cfg.GlobalConfig.Header
		apiTestEnvData.Global = make(map[string]*apistructs.APITestEnvVariable)
		for name, item := range cfg.GlobalConfig.Global {
			apiTestEnvData.Global[name] = &apistructs.APITestEnvVariable{
				Value: item.Value,
				Type:  item.Type,
			}
			caseParams[name] = &apistructs.CaseParams{
				Key:   name,
				Type:  item.Type,
				Value: item.Value,
			}
		}
	}

	// add cookie jar
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if apiTestEnvData != nil && apiTestEnvData.Header != nil && len(apiTestEnvData.Header[CookieJar]) > 0 {
		var cookies cookiejar.Cookies
		err := json.Unmarshal([]byte(apiTestEnvData.Header[CookieJar]), &cookies)
		if err != nil {
			success = false
			logrus.Errorf("failed to unmarshal cookieJar from header, err: %v\n", err)
			return
		}
		cookieJar.SetEntries(cookies)
	}
	hc := http.Client{Jar: cookieJar}
	printGlobalAPIConfig(apiTestEnvData)

	// do apiTest
	apiTest := apitestsv2.New(apiInfo)
	apiReq, apiResp, err := apiTest.Invoke(&hc, apiTestEnvData, caseParams)
	printRenderedHTTPReq(apiReq)
	meta.Req = apiReq
	meta.Resp = apiResp
	meta.CookieJar = cookieJar.GetEntries()
	if apiResp != nil {
		printHTTPResp(apiResp)
	}
	if err != nil {
		meta.Result = resultFailed
		logrus.Errorf("failed to do api test, err: %v", err)
		success = false
		return
	}

	// outParams store in metafile for latter use
	outParams := apiTest.ParseOutParams(apiTest.API.OutParams, apiResp, caseParams)
	printOutParams(outParams, meta)

	// judge asserts
	if len(apiTest.API.Asserts) > 0 {
		// 目前有且只有一组 asserts
		for _, group := range apiTest.API.Asserts {
			succ, assertResults := apiTest.JudgeAsserts(outParams, group)
			printAssertResults(succ, assertResults)
			if !succ {
				addNewLine()
				logrus.Errorf("API Test Success, but asserts failed")
				success = false
				return
			}
		}
	}

	meta.Result = resultSuccess

	addNewLine(2)
	logrus.Println("API Test Success")
}
