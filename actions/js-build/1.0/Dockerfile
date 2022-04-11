FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
# go build
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js-build/1.0/internal/cmd

SHELL ["/bin/bash", "--login", "-c"]

RUN yum -y install gcc bzip2 git

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash

RUN nvm install 12 && nvm install 14 &&  nvm install 16 && nvm install 8 && nvm install 10 && nvm use 12

RUN npm install webpack -g && npm install webpack-cli -g && npm install -g cnpm --registry=https://registry.npmmirror.com
