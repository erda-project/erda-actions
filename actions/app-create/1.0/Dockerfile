FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run actions/app-create/1.0/internal/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-centos:base

RUN yum install -y git
COPY --from=builder /assets /opt/action
