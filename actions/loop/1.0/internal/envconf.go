package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Conf struct {
	Type LoopType `env:"ACTION_TYPE" required:"true"`

	HTTPConf HTTPConf `env:"ACTION_HTTP_CONF"`

	CMD string `env:"ACTION_CMD"`

	JSONResponseSuccessField string `env:"ACTION_JSON_RESPONSE_SUCCESS_FIELD"`

	// 设置循环的最大次数
	LoopMaxTimes uint64 `env:"ACTION_LOOP_MAX_TIMES" default:"10000"`
	// 设置每次循环的间隔时间
	LoopInterval time.Duration `env:"ACTION_LOOP_INTERVAL" default:"10s"`
	// 设置衰退延迟的比例，默认是 1
	LoopDeclineRatio float64 `env:"ACTION_LOOP_DECLINE_RATIO" default:"1"`
	// 设置衰退延迟的最大值，默认不限制最大值
	LoopDeclineLimit time.Duration `env:"ACTION_LOOP_DECLINE_LIMIT" default:"-1s"`
}

func (c *Conf) String() string {
	s := fmt.Sprintf("type: %s\n", c.Type)
	switch c.Type {
	case HTTP:
		b, _ := json.MarshalIndent(c.HTTPConf, "", "  ")
		s += fmt.Sprintf("http_conf: %s\n", string(b))
	case CMD:
		s += fmt.Sprintf("cmd: %s\n", c.CMD)
	}
	if c.JSONResponseSuccessField != "" {
		s += fmt.Sprintf("json_response_success_field: %s\n", c.JSONResponseSuccessField)
	}
	s += fmt.Sprintf(
		"loop_max_times: %d\n"+
			"loop_interval: %fs\n"+
			"loop_decline_ratio: %f\n",
		c.LoopMaxTimes,
		float64(c.LoopInterval)/float64(time.Second),
		c.LoopDeclineRatio,
	)
	if c.LoopDeclineLimit <= 0 {
		s += fmt.Sprintf("loop_decline_limit: no limit\n")
	} else {
		s += fmt.Sprintf("loop_decline_limit: %fs\n", float64(c.LoopDeclineLimit)/float64(time.Second))
	}
	return s
}

type HTTPConf struct {
	Host        string            `json:"host"`
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	QueryParams map[string]string `json:"query_params"`
	Header      map[string]string `json:"header,omitempty"`
	RequestBody string            `json:"request_body"`
}

func (c *Conf) Check() error {

	switch c.Type {
	case HTTP:
		// check http conf
		if c.HTTPConf.Host == "" {
			return errors.Errorf("HTTP loop type missing host, like: 127.0.0.1:8080")
		}
		if c.HTTPConf.Path == "" {
			logger.Println("use default path: /")
			c.HTTPConf.Path = "/"
		}
		switch c.HTTPConf.Method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		case "":
			logger.Println("use default http method: GET")
			c.HTTPConf.Method = http.MethodGet
		default:
			return errors.Errorf("HTTP loop type invalid http method, only support: GET/POST/PUT/DELETE")
		}

	case CMD:
		// check cmd
		if c.CMD == "" {
			return errors.Errorf("CMD loop type missing cmd")
		}

	default:
		return errors.Errorf("invalid loop type, only support: HTTP, CMD")
	}

	return nil
}
