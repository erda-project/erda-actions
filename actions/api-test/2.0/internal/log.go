package main

import (
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
)

func addLineDelimiter(prefix ...string) {
	var _prefix string
	if len(prefix) > 0 {
		_prefix = prefix[0]
	}
	logrus.Printf("%s==========", _prefix)
}

func addNewLine(num ...int) {
	_num := 1
	if len(num) > 0 {
		_num = num[0]
	}
	if _num <= 0 {
		_num = 1
	}
	for i := 0; i < _num; i++ {
		logrus.Println()
	}
}

func printOriginalAPIInfo(api *apistructs.APIInfo) {
	if api == nil {
		return
	}
	logrus.Printf("Original API Info:")
	defer addNewLine()
	// name
	if api.Name != "" {
		logrus.Printf("name: %s", api.Name)
	}
	// url
	logrus.Printf("url: %s", api.URL)
	// method
	logrus.Printf("method: %s", api.Method)
	// headers
	if len(api.Headers) > 0 {
		logrus.Printf("headers:")
		for _, h := range api.Headers {
			logrus.Printf("  key: %s", h.Key)
			logrus.Printf("  value: %s", h.Value)
			if h.Desc != "" {
				logrus.Printf("  desc: %s", h.Desc)
			}
			addLineDelimiter("  ")
		}
	}
	// params
	if len(api.Params) > 0 {
		logrus.Printf("params:")
		for _, p := range api.Params {
			logrus.Printf("  key: %s", p.Key)
			logrus.Printf("  value: %s", p.Value)
			if p.Desc != "" {
				logrus.Printf("  desc: %s", p.Desc)
			}
			addLineDelimiter("  ")
		}
	}
	// request body
	if api.Body.Type != "" {
		logrus.Printf("request body:")
		logrus.Printf("  type: %s", api.Body.Type.String())
		logrus.Printf("  content: %s", jsonOneLine(api.Body.Content))
	}
	// out params
	if len(api.OutParams) > 0 {
		logrus.Printf("out params:")
		for _, out := range api.OutParams {
			logrus.Printf("  arg: %s", out.Key)
			logrus.Printf("  source: %s", out.Source.String())
			if out.Expression != "" {
				logrus.Printf("  expr: %s", out.Expression)
			}
			addLineDelimiter("  ")
		}
	}
	// asserts
	if len(api.Asserts) > 0 {
		logrus.Printf("asserts:")
		for _, group := range api.Asserts {
			for _, assert := range group {
				logrus.Printf("  key: %s", assert.Arg)
				logrus.Printf("  operator: %s", assert.Operator)
				logrus.Printf("  value: %s", assert.Value)
				addLineDelimiter("  ")
			}
		}
	}
}

func printGlobalAPIConfig(cfg *apistructs.APITestEnvData) {
	if cfg == nil {
		return
	}
	logrus.Printf("Global API Config:")
	defer addNewLine()

	// name
	if cfg.Name != "" {
		logrus.Printf("name: %s", cfg.Name)
	}
	// domain
	logrus.Printf("domain: %s", cfg.Domain)
	// headers
	if len(cfg.Header) > 0 {
		logrus.Printf("headers:")
		for k, v := range cfg.Header {
			logrus.Printf("  key: %s", k)
			logrus.Printf("  value: %s", v)
			addLineDelimiter("  ")
		}
	}
	// global
	if len(cfg.Global) > 0 {
		logrus.Printf("global configs:")
		for key, item := range cfg.Global {
			logrus.Printf("  key: %s", key)
			logrus.Printf("  value: %s", item.Value)
			logrus.Printf("  type: %s", item.Type)
			if item.Desc != "" {
				logrus.Printf("  desc: %s", item.Desc)
			}
			addLineDelimiter("  ")
		}
	}
}

func printRenderedHTTPReq(req *apistructs.APIRequestInfo) {
	if req == nil {
		return
	}
	logrus.Printf("Rendered HTTP Request:")
	defer addNewLine()

	// url
	logrus.Printf("url: %s", req.URL)
	// method
	logrus.Printf("method: %s", req.Method)
	// headers
	if len(req.Headers) > 0 {
		logrus.Printf("headers:")
		for key, values := range req.Headers {
			logrus.Printf("  key: %s", key)
			if len(values) == 1 {
				logrus.Printf("  value: %s", values[0])
			} else {
				logrus.Printf("  values: %v", values)
			}
			addLineDelimiter("  ")
		}
	}
	// params
	if len(req.Params) > 0 {
		logrus.Printf("params:")
		for key, values := range req.Params {
			logrus.Printf("  key: %s", key)
			if len(values) == 1 {
				logrus.Printf("  value: %s", values[0])
			} else {
				logrus.Printf("  values: %v", values)
			}
			addLineDelimiter("  ")
		}
	}
	// body
	if req.Body.Type != "" {
		logrus.Printf("request body:")
		logrus.Printf("  type: %s", req.Body.Type.String())
		logrus.Printf("  content: %s", req.Body.Content)
	}
}

func printHTTPResp(resp *apistructs.APIResp) {
	if resp == nil {
		return
	}
	logrus.Printf("HTTP Response:")
	defer addNewLine()

	// status
	logrus.Printf("http status: %d", resp.Status)
	// headers
	if len(resp.Headers) > 0 {
		logrus.Printf("response headers:")
		for key, values := range resp.Headers {
			logrus.Printf("  key: %s", key)
			if len(values) == 1 {
				logrus.Printf("  value: %s", values[0])
			} else {
				logrus.Printf("  values: %v", values)
			}
			addLineDelimiter("  ")
		}
	}
	// response body
	if resp.BodyStr != "" {
		logrus.Printf("response body: %s", resp.BodyStr)
	}
}

func printOutParams(outParams map[string]interface{}, meta *Meta) {
	if len(outParams) == 0 {
		return
	}
	logrus.Printf("Out Params:")
	defer addNewLine()

	// 按定义顺序返回
	for _, define := range meta.OutParamsDefine {
		k := define.Key
		v, ok := outParams[k]
		if !ok {
			continue
		}
		meta.OutParamsResult[k] = v
		logrus.Printf("  arg: %s", k)
		logrus.Printf("  source: %s", define.Source.String())
		if define.Expression != "" {
			logrus.Printf("  expr: %s", define.Expression)
		}
		logrus.Printf("  value: %s", jsonOneLine(v))
		var vtype string
		if v == nil {
			vtype = "nil"
		} else {
			vtype = reflect.TypeOf(v).String()
		}
		logrus.Printf("  type: %s", vtype)
		addLineDelimiter("  ")
	}
}

func printAssertResults(success bool, results []*apistructs.APITestsAssertData) {
	logrus.Printf("Assert Result: %t", success)
	defer addNewLine()

	logrus.Printf("Assert Detail:")
	for _, result := range results {
		logrus.Printf("  arg: %s", result.Arg)
		logrus.Printf("  operator: %s", result.Operator)
		logrus.Printf("  value: %s", result.Value)
		logrus.Printf("  actualValue: %s", jsonOneLine(result.ActualValue))
		logrus.Printf("  success: %t", result.Success)
		if result.ErrorInfo != "" {
			logrus.Printf("  errorInfo: %s", result.ErrorInfo)
		}
		addLineDelimiter("  ")
	}
}
