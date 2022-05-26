# Copyright (c) 2021 Terminus, Inc.
#
# This program is free software: you can use, redistribute, and/or modify
# it under the terms of the GNU Affero General Public License, version 3
# or later ("AGPL"), as published by the Free Software Foundation.
#
# This program is distributed in the hope that it will be useful, but WITHOUT
# ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
# FITNESS FOR A PARTICULAR PURPOSE.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.

FROM --platform=linux/amd64 registry.erda.cloud/erda/centos:7

MAINTAINER chenzhongrun "zhongrun.czr@alibaba-inc.com"

# build go
RUN wget https://golang.google.cn/dl/go1.17.7.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz

# build cmd
COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions/actions
RUN rm -rf `ls | grep -v project-artifacts`
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN /usr/local/go/bin/go env -w GOPROXY="https://goproxy.cn,direct"
RUN /usr/local/go/bin/go env -w GOOS="linux"
RUN /usr/local/go/bin/go env -w GOARCH="amd64"
RUN /usr/local/go/bin/go mod tidy
RUN /usr/local/go/bin/go build -o /opt/action/run actions/project-artifacts/1.0/internal/cmd/main.go

RUN chmod 777 /opt/action/*
