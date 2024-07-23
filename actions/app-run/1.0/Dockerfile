FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/golang:1.22-bookworm AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
ENV GO111MODULE=on
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run actions/app-run/1.0/internal/cmd/*.go
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/when_sigterm  actions/app-run/1.0/internal/appCancel/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-centos:base
COPY --from=builder /assets /opt/action
