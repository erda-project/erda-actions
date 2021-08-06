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


## action 编译
FROM golang:1.15 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ENV GOPROXY https://goproxy.cn/,direct
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run actions/erda-pkg-release-public/1.0/internal/cmd/main.go


## work image
FROM centos

## maintainer
MAINTAINER xujie "xj236673@alibaba-inc.com"

## linux tools install
RUN alias ll='ls -l' && \
    yum install -y yum-utils tree git wget make vim make && \
    yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo && \
    yum -y install docker-ce-cli

## oss client prepare
RUN wget http://gosspublic.alicdn.com/ossutil/1.7.3/ossutil64 && \
    chmod 755 ossutil64 && \
    cp ossutil64 /usr/bin/

## golang prepare
ENV GOLANG_VERSION 1.16.4
RUN set -eux; \
	\
	goRelArch='linux-amd64'; goRelSha256='7154e88f5a8047aad4b80ebace58a059e36e7e2e4eb3b383127a28c711b4ff59'; \
	\
	url="https://golang.google.cn/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
	curl -L -o go.tgz "$url"; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz

ENV GO111MODULE on
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV GOPROXY https://goproxy.cn/,direct

## action caller entry
COPY --from=builder /opt/action/run /opt/action/run
RUN chmod 777 /opt/action/*
