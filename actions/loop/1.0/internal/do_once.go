package main

import (
	"bytes"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thedevsaddam/gojsonq/v2"

	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
)

func (o *Object) doOnceHTTP() (bool, error) {

	var request *httpclient.Request

	switch o.conf.HTTPConf.Method {
	case http.MethodPost:
		request = o.hc.Post(o.conf.HTTPConf.Host)
	case http.MethodPut:
		request = o.hc.Put(o.conf.HTTPConf.Host)
	case http.MethodDelete:
		request = o.hc.Delete(o.conf.HTTPConf.Host)
	default:
		request = o.hc.Get(o.conf.HTTPConf.Host, httpclient.RetryOption{})
	}

	request = request.Path(o.conf.HTTPConf.Path)
	// query params
	for k, v := range o.conf.HTTPConf.QueryParams {
		request = request.Param(k, v)
	}
	// header
	for k, v := range o.conf.HTTPConf.Header {
		request = request.Header(k, v)
	}
	// json body
	if len(o.conf.HTTPConf.RequestBody) > 0 {
		request = request.Header("Content-Type", "application/json").RawBody(bytes.NewBufferString(o.conf.HTTPConf.RequestBody))
	}

	// response
	var body bytes.Buffer
	resp, err := request.Do().Body(&body)
	if err != nil {
		return false, errors.Errorf("failed to do http request, err: %v", err)
	}
	bodyStr := body.String()
	logger.Printf("response body: %s\n", bodyStr)
	if !resp.IsOK() {
		return false, errors.Errorf("http code is not 2xx, http code: %d", resp.StatusCode())
	}
	// use special success json field
	if o.conf.JSONResponseSuccessField != "" {
		accessPath := strutil.TrimLeft(o.conf.JSONResponseSuccessField, ".")
		v := gojsonq.New().FromString(bodyStr).Find(accessPath)
		if v != nil {
			success, _ := v.(bool)
			if success {
				return success, nil
			}
		}
		logger.Printf("%q from JSON Response Body is not `true`", o.conf.JSONResponseSuccessField)
		return false, nil
	}
	return true, nil
}

func (o *Object) doOnceCMD() (bool, error) {
	return false, errors.Errorf("cmd not support now")
}
