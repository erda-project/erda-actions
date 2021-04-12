package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda/apistructs"
)

func makeConfFromEnv(envConfig EnvConfig) {
	// url
	_ = os.Setenv("ACTION_URL", envConfig.URL)
	// method
	_ = os.Setenv("ACTION_METHOD", envConfig.Method)
	// params
	paramsByte, err := json.Marshal(envConfig.Params)
	if err != nil {
		logrus.Fatalf("invalid params, err: %v", err)
	}
	_ = os.Setenv("ACTION_PARAMS", string(paramsByte))
	// headers
	headersByte, err := json.Marshal(envConfig.Headers)
	if err != nil {
		logrus.Fatalf("invalid headers, err: %v", err)
	}
	_ = os.Setenv("ACTION_HEADERS", string(headersByte))
	// body
	bodyByte, err := json.Marshal(envConfig.Body)
	if err != nil {
		logrus.Fatalf("invalid body, err: %v", err)
	}
	_ = os.Setenv("ACTION_BODY", string(bodyByte))
	// out params
	outParamsByte, err := json.Marshal(envConfig.OutParams)
	if err != nil {
		logrus.Fatalf("invalid outParams, err: %v", err)
	}
	_ = os.Setenv("ACTION_OUT_PARAMS", string(outParamsByte))
	// asserts
	assertsByte, err := json.Marshal(envConfig.Asserts)
	if err != nil {
		logrus.Fatalf("invalid asserts, err: %v", err)
	}
	_ = os.Setenv("ACTION_ASSERTS", string(assertsByte))
	// global config
	globalConfigByte, err := json.Marshal(envConfig.GlobalConfig)
	if err != nil {
		logrus.Fatalf("invalid globalConfig, err: %v", err)
	}
	_ = os.Setenv("AUTOTEST_API_GLOBAL_CONFIG", string(globalConfigByte))
}

func TestMain_Addons(t *testing.T) {
	makeConfFromEnv(EnvConfig{
		URL:    "/api/addons",
		Method: "GET",
		Params: []APIParam{
			{
				Key:   "type",
				Value: "project",
			},
			{
				Key:   "value",
				Value: 2,
			},
		},
		Headers: []APIHeader{
			{
				Key:   "h1",
				Value: "v1",
				Desc:  "",
			},
		},
		OutParams: []apistructs.APIOutParam{
			{
				Key:        "list",
				Source:     apistructs.APIOutParamSourceBodyJson,
				Expression: ".data",
			},
			{
				Key:        "monitorName",
				Source:     apistructs.APIOutParamSourceBodyJson,
				Expression: ".data[0].addonName",
			},
			{
				Key:        "monitorAttachCount",
				Source:     apistructs.APIOutParamSourceBodyJson,
				Expression: ".data[0].attachCount",
			},
			{
				Key:        "monitorCanDel",
				Source:     apistructs.APIOutParamSourceBodyJson,
				Expression: ".data[0].canDel",
			},
			{
				Key:        "monitorConfig",
				Source:     apistructs.APIOutParamSourceBodyJson,
				Expression: ".data[0].config",
			},
			{
				Key:        "monitorConfigPublicHost",
				Source:     apistructs.APIOutParamSourceBodyJson,
				Expression: ".data[0].config.PUBLIC_HOST",
			},
		},
		Asserts: []APIAssert{
			{
				Arg:      "list",
				Operator: "contains",
				Value:    `instanceId`,
			},
			{
				Arg:      "monitorName",
				Operator: "=",
				Value:    "monitor",
			},
			{
				Arg:      "monitorAttachCount",
				Operator: "=",
				Value:    "3",
			},
			{
				Arg:      "monitorCanDel",
				Operator: "=",
				Value:    "false",
			},
			{
				Arg:      "monitorConfig",
				Operator: "contains",
				Value:    "TERMINUS_AGENT_ENABLE",
			},
			{
				Arg:      "monitorConfigPublicHost",
				Operator: "=",
				Value:    "",
			},
		},
		GlobalConfig: &apistructs.AutoTestAPIConfig{
			Domain: "https://terminus-test-org.test.terminus.io",
			Header: map[string]string{
				"Cookie": "OPENAPISESSION=b670260f-85c6-40a4-9cdc-f35c71c84722",
			},
			Global: nil,
		},
		MetaFile: "",
	})
	main()
}

func TestMain_CreatIssue(t *testing.T) {
	makeConfFromEnv(EnvConfig{
		URL:     "/api/issues",
		Method:  http.MethodPost,
		Headers: nil,
		Body: apistructs.APIBody{
			Type: apistructs.APIBodyTypeApplicationJSON,
			Content: `
{
  "projectID":2,
  "iterationID":8,
  "priority":"NORMAL",
  "complexity":"NORMAL",
  "severity":"NORMAL",
  "taskType":"dev",
  "bugStage":"codeDevelopment",
  "title":"2",
  "assignee":"2",
  "planStartedAt":"2020-12-04T00:00:00+08:00",
  "planFinishedAt":"2020-12-04T00:00:00+08:00",
  "issueManHour":{
    "estimateTime":180,
    "remainingTime":180
  },
  "type":"TASK"
}
`,
		},
		OutParams: []apistructs.APIOutParam{
			{
				Key:    "status",
				Source: apistructs.APIOutParamSourceStatus,
			},
		},
		Asserts: []APIAssert{
			{
				Arg:      "status",
				Operator: "=",
				Value:    "200",
			},
		},
		GlobalConfig: &apistructs.AutoTestAPIConfig{
			Domain: "https://terminus-test-org.test.terminus.io",
			Header: map[string]string{
				"Cookie": "OPENAPISESSION=b670260f-85c6-40a4-9cdc-f35c71c84722",
			},
			Global: nil,
		},
		MetaFile: "",
	})
	main()
}

func TestMain_UrlEncoded(t *testing.T) {
	makeConfFromEnv(EnvConfig{
		URL:     "/test",
		Method:  http.MethodPost,
		Headers: nil,
		Body: apistructs.APIBody{
			Type: apistructs.APIBodyTypeApplicationXWWWFormUrlencoded,
			Content: []apistructs.APIParam{
				{
					Key:   "userID",
					Value: "{{userID}}",
				},
				{
					Key:   "name",
					Value: "{{name}}",
					Desc:  "{{name}} desc",
				},
			},
		},
		GlobalConfig: &apistructs.AutoTestAPIConfig{
			Domain: "https://terminus-test-org.test.terminus.io",
			Header: map[string]string{
				"Cookie": "OPENAPISESSION=b670260f-85c6-40a4-9cdc-f35c71c84722",
			},
			Global: map[string]apistructs.AutoTestConfigItem{
				"userID": {Type: "string", Value: "2"},
				"name":   {Type: "string", Value: "my-name"},
			},
		},
		OutParams: []apistructs.APIOutParam{
			{
				Key:    "status",
				Source: apistructs.APIOutParamSourceStatus,
			},
		},
		Asserts: []APIAssert{
			{
				Arg:      "status",
				Operator: "!=",
				Value:    "200",
			},
		},
	})
	main()
}

func TestMain_CreateManualTestCase(t *testing.T) {
	makeConfFromEnv(EnvConfig{
		URL:     "/api/testcases",
		Method:  http.MethodPost,
		Headers: nil,
		Body: apistructs.APIBody{
			Type: apistructs.APIBodyTypeApplicationXWWWFormUrlencoded,
			Content: []apistructs.APIParam{
				{
					Key:   "userID",
					Value: "{{userID}}",
				},
				{
					Key:   "name",
					Value: "{{name}}",
					Desc:  "{{name}} desc",
				},
			},
		},
		GlobalConfig: &apistructs.AutoTestAPIConfig{
			Domain: "https://terminus-test-org.test.terminus.io",
			Header: map[string]string{
				"Cookie": "OPENAPISESSION=xxx",
			},
			Global: map[string]apistructs.AutoTestConfigItem{
				"userID": {Type: "string", Value: "2"},
				"name":   {Type: "string", Value: "my-name"},
			},
		},
		OutParams: []apistructs.APIOutParam{
			{
				Key:    "status",
				Source: apistructs.APIOutParamSourceStatus,
			},
		},
		Asserts: []APIAssert{
			{
				Arg:      "status",
				Operator: "!=",
				Value:    "200",
			},
		},
	})
	main()
}

func TestMain_GaiaProduct(t *testing.T) {
	makeConfFromEnv(EnvConfig{
		URL:    "http://test-gateway.app.terminus.io/t-product/itemcenter-backend/itemcenter/10176/api/trantor/action/exe1",
		Method: http.MethodPost,
		Headers: []APIHeader{
			{
				Key:   "Content-Type",
				Value: "application/json",
			},
			{
				Key:   "Cookie",
				Value: "lng=zh-CN; t_product_test_u_c_local=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJsb2dpbiIsInBhdGgiOiIvIiwidG9rZW5LZXkiOiIzNTU4OTZhNWUxN2E1NjNjMTI2ZThkMGY0ODhkZjgyZTllNmFhZjUwYWQ3MmJkN2RiY2Q3NDAxOTFiZjMxNGJhIiwibmJmIjoxNjExMDI3ODQ2LCJkb21haW4iOiJ0ZXJtaW51cy5pbyIsImlzcyI6ImRyYWNvIiwidGVuYW50SWQiOjEsImV4cGlyZV90aW1lIjo2MDQ4MDAsImV4cCI6MTYxMTYzMjY0NiwiaWF0IjoxNjExMDI3ODQ2fQ.hCo4xKxrNBjZV0muSTAJgshyB8i8zt5DH4s68Hh_rxs",
			},
		},
		Body: apistructs.APIBody{
			Type: apistructs.APIBodyTypeApplicationJSON,
			Content: `
{
  "frontendContext":{

  },
  "actionKey":"itemcenter_ItemVO_ItemAction::create",
  "context":{
      "modelKey":"itemcenter_ItemVO",
      "actionLabel":"提交",
      "record":[
          {
              "isCombined":false,
              "category":{
                  "id":34008
              },
              "type":1,
              "name":"测试55",
              "spu":null,
              "brand":null,
              "advertise":null,
              "mainImageAttachment":{
                  "files":[
                      {
                          "name":"20200909154255.jpg",
                          "url":"//terminus-trantor.oss-cn-hangzhou.aliyuncs.com/trantor/attachments/b2c8dc12-3ca3-4203-a57e-655ea7fbe243.jpg",
                          "type":"jpg",
                          "size":50717
                      }
                  ]
              },
              "videoUrl":null,
              "taxRate":"5%",
              "unit":null,
              "version":null,
              "keyword":null,
              "deliveryFeeTemplates":{
                  "middleFee":0,
                  "fee":0,
                  "lowFee":0,
                  "sellerId":10083,
                  "isFree":true,
                  "highFee":0,
                  "lowPrice":0,
                  "isSpecial":true,
                  "highPrice":0,
                  "incrFee":0,
                  "id":8001,
                  "updatedAt":1606135255000,
                  "isDefault":true,
                  "initFee":0,
                  "chargeMethod":2,
                  "name":"免运费",
                  "status":1
              },
              "categoryAttributes":[

              ],
              "skuAttributes":[

              ],
              "skuList":[
                  {
                      "id":1611026682967,
                      "enable":true,
                      "barcode":"1",
                      "originalPrice":1000,
                      "price":1000,
                      "unitAmount":"1"
                  }
              ],
              "otherAttributes":[
                  {
                      "group":"BASIC",
                      "otherAttributes":[

                      ]
                  },
                  {
                      "group":"DEFAULT",
                      "otherAttributes":[

                      ]
                  },
                  {
                      "group":"USER_DEFINED",
                      "otherAttributes":[

                      ]
                  }
              ]
          }
      ]
  }
}`,
		},
		GlobalConfig: &apistructs.AutoTestAPIConfig{
			Domain: "http://test-gateway.app.terminus.io",
		},
		//OutParams: []apistructs.APIOutParam{
		//	{
		//		Key:    "status",
		//		Source: apistructs.APIOutParamSourceStatus,
		//	},
		//},
		//Asserts: []APIAssert{
		//	{
		//		Arg:      "status",
		//		Operator: "!=",
		//		Value:    "200",
		//	},
		//},
	})
	main()
}

func TestMain_GaiaProductOrderId(t *testing.T) {
	os.Setenv("ACTION_PARAMS", `[{"key":"orderId","value":1352141084883972097}]`)
	os.Setenv("ACTION_URL", "x")
	os.Setenv("ACTION_METHOD", "GET")
	main()
}

func TestMain_AcceptEncoding(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://bbc-distribution-test.app.terminus.io/api/distribution/config/info", bytes.NewBufferString(`{"type":"application/json"}`))
	assert.NoError(t, err)
	req.Header.Add("Accept-Encoding", "identity")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	spew.Dump(resp.TransferEncoding)
	rr, err := gzip.NewReader(resp.Body)
	assert.NoError(t, err)
	bodyByte, _ := ioutil.ReadAll(rr)
	body2 := string(bodyByte)
	defer resp.Body.Close()
	fmt.Println(body2)
}
