package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/api-test/2.0/internal/cookiejar"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	resultSuccess = "success"
	resultFailed  = "failed"
)

const (
	metaKeyResult           = "result"
	metaKeyAPIRequest       = "api_request"
	metaKeyAPIResponse      = "api_response"
	metaKeyAPICookies       = "api_set_cookies"
	metaKeyAPIAssertSuccess = "api_assert_success" // true; false
	metaKeyAPIAssertDetail  = "api_assert_detail"
)

type Meta struct {
	Result          string
	AssertResult    bool
	AssertDetail    string
	Req             *apistructs.APIRequestInfo
	Resp            *apistructs.APIResp
	OutParamsDefine []apistructs.APIOutParam
	CookieJar       cookiejar.Cookies
	OutParamsResult map[string]interface{}
}

func NewMeta() *Meta {
	return &Meta{
		OutParamsResult: map[string]interface{}{},
	}
}

type KVs []kv
type kv struct {
	k string
	v string
}

func (kvs *KVs) add(k, v string) {
	*kvs = append(*kvs, kv{k, v})
}

func writeMetaFile(metafilePath string, meta *Meta) {
	var content string

	// kvs 保证顺序
	kvs := &KVs{}

	kvs.add(metaKeyResult, meta.Result)
	if meta.AssertDetail != "" {
		kvs.add(metaKeyAPIAssertSuccess, strconv.FormatBool(meta.AssertResult))
		kvs.add(metaKeyAPIAssertDetail, meta.AssertDetail)
	}
	if meta.Req != nil {
		kvs.add(metaKeyAPIRequest, jsonOneLine(meta.Req))
	}
	if meta.Resp != nil {
		kvs.add(metaKeyAPIResponse, jsonOneLine(meta.Resp))
	}
	if meta.CookieJar != nil {
		if meta.Resp != nil && meta.Resp.Headers != nil && meta.Resp.Headers["Set-Cookie"] != nil {
			jar, _ := json.Marshal(meta.CookieJar)
			kvs.add(metaKeyAPICookies, string(jar))
		}
	}
	if len(meta.OutParamsResult) > 0 {
		for _, define := range meta.OutParamsDefine {
			v, ok := meta.OutParamsResult[define.Key]
			if !ok {
				continue
			}
			kvs.add(define.Key, jsonOneLine(v))
		}
	}

	for _, kv := range *kvs {
		content = fmt.Sprintf("%s\n%s=%s\n", content, kv.k, kv.v)
	}
	if err := filehelper.CreateFile(metafilePath, content, 0644); err != nil {
		logrus.Printf("failed to create metafile, err: %v\n", err)
	}
}
