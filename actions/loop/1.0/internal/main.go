package main

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type LoopType string

const (
	HTTP = LoopType("HTTP")
	CMD  = LoopType("CMD")
)

var logger = log.New(os.Stdout, "[Loop Action] ", 0)

type Object struct {
	conf *Conf

	hc *httpclient.HTTPClient

	successFlag *bool
}

func Initialize() (*Object, error) {

	// get conf from env
	var conf Conf
	if err := envconf.Load(&conf); err != nil {
		return nil, errors.Errorf("failed to get conf from env, err: %v", err)
	}
	if err := conf.Check(); err != nil {
		return nil, err
	}
	// print action params
	logger.Printf("Loop Config:\n%s\n", conf.String())

	// http client
	hc := httpclient.New(
		httpclient.WithCompleteRedirect(),
		httpclient.WithTimeout(time.Second, time.Second*3),
	)

	obj := Object{
		conf:        &conf,
		hc:          hc,
		successFlag: &[]bool{false}[0],
	}

	return &obj, nil
}

func main() {
	obj, err := Initialize()
	if err != nil {
		logger.Fatalf("failed to initialize action, err: %v\n", err)
	}

	if err := Loop(obj); err != nil {
		logger.Fatalf("loop failed, err: %v\n", err)
	}
}
