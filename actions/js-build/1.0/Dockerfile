FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG TARGETARCH

# go build
RUN GOOS=linux GOARCH=$TARGETARCH go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js-build/1.0/internal/cmd

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/oraclelinux:7

COPY --from=builder /opt/action/run /opt/action/run

RUN chmod 777 /opt/action/*

SHELL ["/bin/sh", "--login", "-c"]

RUN yum -y install gcc bzip2 git

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash

RUN source ~/.bashrc && nvm install 12 && nvm install 14 && nvm install 16 && nvm install 8 && nvm install 10 && nvm use 12

RUN source ~/.bashrc && npm install webpack -g && npm install webpack-cli -g && npm install -g cnpm --registry=https://registry.npmmirror.com
