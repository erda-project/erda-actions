package apitest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"

	"github.com/erda-project/erda-actions/actions/api-test/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/apitestsv2"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/http/httpserver/errorresp"
)

func handleAPIs(apiIDs []uint64) error {
	var (
		err        error
		caseParams = make(map[string]*apistructs.CaseParams)
		envInfo    = &apistructs.APITestEnvData{
			Header: make(map[string]string),
			Global: make(map[string]*apistructs.APITestEnvVariable),
		}
	)

	envInfo, err = getProjectEnvInfo()
	if err != nil {
		logrus.Warningf("failed to get project test env info, (%+v)", err)
	}

	// get usecase env variable
	usecaseID := strconv.FormatUint(conf.UsecaseID(), 10)
	if usecaseID == "" {
		return errors.Errorf("empty usecase ID")
	}

	usecaseEnvData, err := getUsecaseTestEnvInfo(usecaseID)
	if err != nil {
		logrus.Warningf("not exist usecase test env info, usecaseID:%s, (%+v)", usecaseID, err)
	}

	if usecaseEnvData != nil {
		for k, v := range usecaseEnvData.Global {
			envInfo.Global[k] = v
		}

		for k, v := range usecaseEnvData.Header {
			envInfo.Header[k] = v
		}
	}

	// render project env global params, least low priority
	if envInfo != nil && envInfo.Global != nil {
		for k, v := range envInfo.Global {
			caseParams[k] = &apistructs.CaseParams{
				Type:  v.Type,
				Value: v.Value,
			}
		}
	}

	// add cookie jar
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		logrus.Warningf("failed to new cookie jar")
	}

	httpClient := &http.Client{}

	if cookieJar != nil {
		httpClient.Jar = cookieJar
	}

	for _, apiID := range apiIDs {
		// 单个 API 执行失败，不返回失败，继续执行下一个
		if err := handleOneAPI(httpClient, envInfo, apiID, caseParams); err != nil {
			logrus.Warningf("handle api error, apiID:%d, (%+v)", apiID, err)
		}
	}

	return nil
}

func handleOneAPI(httpClient *http.Client, envInfo *apistructs.APITestEnvData, apiID uint64, caseParams map[string]*apistructs.CaseParams) error {
	// 根据 apiID 查询 api
	apiTestInfo, err := getAPITestInfo(apiID)
	if err != nil {
		return err
	}
	// 转换为 apistructs.APIInfo
	var apiInfo apistructs.APIInfo
	if err := json.Unmarshal([]byte(apiTestInfo.ApiInfo), &apiInfo); err != nil {
		return err
	}
	apiInfo.ID = strconv.FormatUint(apiID, 10)

	apiTest := apitestsv2.New(&apiInfo, apitestsv2.WithTryV1RenderJsonBodyFirst())

	logrus.Infof("<<<<< Start execute API(%s) ...>>>>>", apiTest.API.URL)
	apiReq, resp, err := apiTest.Invoke(httpClient, envInfo, caseParams)
	if err != nil {
		// 失败，打出警告
		logrus.Warningf("invoke api error, apiInfo:%+v, (%+v)", apiTest, err)

		// 上报错误
		respData := &apistructs.APIResp{
			BodyStr: err.Error(),
		}
		respStr, err := json.Marshal(respData)
		if err != nil {
			respStr = []byte(fmt.Sprint(respStr))
		}

		reqStr, err := json.Marshal(apiReq)
		if err != nil {
			reqStr = []byte(fmt.Sprint(reqStr))
		}

		reportResult(apiID, false, string(reqStr), string(respStr), "")
		return err
	}
	logrus.Infof("invoke status code: %d, \nheaders: %v, \nbody: %v", resp.Status, resp.Headers, string(resp.Body))

	outParams := apiTest.ParseOutParams(apiTest.API.OutParams, resp, caseParams)

	assertResult := &apistructs.APITestsAssertResult{}
	if len(apiTest.API.Asserts) > 0 {
		succ, assertResultData := apiTest.JudgeAsserts(outParams, apiTest.API.Asserts[0])
		logrus.Infof("judge assert result: %v", succ)
		assertResult.Success = succ
		assertResult.Result = assertResultData
	}

	respStr, err := json.Marshal(resp)
	if err != nil {
		respStr = []byte(fmt.Sprint(resp))
	}

	reqStr, err := json.Marshal(apiReq)
	if err != nil {
		reqStr = []byte(fmt.Sprint(reqStr))
	}

	assertStr, err := json.Marshal(assertResult)
	if err != nil {
		assertStr = []byte(fmt.Sprint(assertResult))
	}

	return reportResult(apiID, assertResult.Success, string(reqStr), string(respStr), string(assertStr))
}

func getAPITestInfo(apiID uint64) (*apistructs.ApiTestInfo, error) {
	var resp apistructs.ApiTestsGetResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/apitests/%d", apiID)).
		Header("Authorization", conf.DiceOpenapiToken()).
		Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() || !resp.Success {
		return nil, errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}
	return resp.Data, nil
}

func reportResult(apiID uint64, succ bool, apiReq, apiResponse, assertResult string) error {
	var result apistructs.ApiTestInfo
	if succ {
		result.Status = apistructs.ApiTestPassed
	} else {
		result.Status = apistructs.ApiTestFailed
	}
	result.AssertResult = assertResult
	result.ApiRequest = apiReq
	result.ApiResponse = apiResponse

	// invoke
	params := make(url.Values)
	params.Add("isResult", "true")
	var resp apistructs.ApiTestsUpdateResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Put(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/apitests/%d", apiID)).
		Header("Authorization", conf.DiceOpenapiToken()).
		Params(params).
		JSONBody(result).Do().JSON(&resp)
	if err != nil {
		return err
	}
	if !r.IsOK() || !resp.Success {
		return errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}
	return nil
}

func getProjectEnvInfo() (*apistructs.APITestEnvData, error) {
	// invoke
	var resp apistructs.APITestEnvGetResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/testenv/%d", conf.ProjectTestEnvID())).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() || !resp.Success {
		return nil, errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}

	return resp.Data, nil
}

func getUsecaseTestEnvInfo(envID string) (*apistructs.APITestEnvData, error) {
	// invoke
	var resp apistructs.APITestEnvListResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path("/api/testenv/actions/list-envs").
		Param("envID", envID).
		Param("envType", "case").
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() || !resp.Success {
		return nil, errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}

	if len(resp.Data) > 0 {
		return resp.Data[0], nil
	}

	return nil, nil
}
