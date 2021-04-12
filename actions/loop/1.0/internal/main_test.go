package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/thedevsaddam/gojsonq/v2"
)

func TestEnvConf(t *testing.T) {
	httpConf := HTTPConf{
		Host:        "http://dev-api-gateway.kube-system.svc.cluster.local",
		Path:        "/exxonmobil-member-cdp/api/cdp/workflows/actions/check-dependencies",
		Method:      http.MethodPost,
		QueryParams: nil,
		RequestBody: `{"workflowId":183,"dependWorkflowIds":[181,182]}`,
	}
	b, _ := json.Marshal(&httpConf)
	fmt.Println(string(b))
}

func TestJSONFieldAccess(t *testing.T) {
	s := `[{"host":"localhost:3081","path":"/ping","method":"","queryParams":{"name":"linjun"},"requestBody":null},{"host":"localhost:3081","path":"/ping","method":"","queryParams":{"name":"linjun"},"requestBody":null}]`
	v := gojsonq.New().FromString(s).Find("[0].queryParams.name")
	fmt.Printf("%v\n", v)
}
