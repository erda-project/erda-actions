FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
# go build
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js-build/1.0/internal/cmd

SHELL ["/bin/bash", "--login", "-c"]

RUN yum -y install gcc bzip2 git

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash

RUN nvm install v12.18.4 && nvm install v14.11.0 &&  nvm install v8.17.0 && nvm install v10.22.1 &&  nvm use v12.18.4

RUN npm install webpack -g && npm install webpack-cli -g && npm install -g cnpm --registry=https://registry.npm.taobao.org
