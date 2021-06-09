// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/push-extensions/1.0/internal/client"
	"github.com/erda-project/erda-actions/actions/push-extensions/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/push-extensions/1.0/internal/workdir"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

var (
	repos    = flag.String("repos", ".", "repos")
	username = flag.String("username", "", "username")
	password = flag.String("password", "", "password")
	openhost = flag.String("host", "", "openapi host")
)

func main() {
	flag.Parse()

	// load configuration and pre-check
	conf := new(config.Config)
	conf.Repos = strings.Split(*repos, ",")
	conf.Username = *username
	conf.Password = *password
	conf.Host = *openhost

	if len(conf.Repos) == 0 {
		err := errors.New("there is no repo, did you add it ?")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}
	if conf.Host == "" || conf.Username == "" || conf.Password == "" {
		err := errors.New("missing parameters, did you set the openapi host, username and password ?")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	// make extension pushing client
	cli, err := client.New(conf.Host, conf.Username, conf.Password)
	if err != nil {
		err = errors.Wrap(err, "failed to client.New")
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalln(err)
	}

	// make context
	ctx := context.Background()
	ctx = context.WithValue(ctx, workdir.CtxKeyConfig, conf)
	ctx = context.WithValue(ctx, workdir.CtxKeyClient, cli)

	// load and push all extensions from every repo
	for _, repo := range conf.Repos {
		extensionsRepo := workdir.New(ctx, repo)

		if err := extensionsRepo.LoadExtensions(); err != nil {
			err = errors.Wrapf(err, "failed to LoadExtensions from %s", repo)
			_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
			logrus.Fatalln(err)
		}

		if err := extensionsRepo.Push(); err != nil {
			err = errors.Wrapf(err, "failed to Push extensions from %s", repo)
			_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
			logrus.Fatalln(err)
		}
	}

	_ = metawriter.Write(map[string]interface{}{config.Success: true})
	os.Exit(0)
}
