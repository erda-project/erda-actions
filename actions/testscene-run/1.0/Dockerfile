FROM --platform=$TARGETPLATFORM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run actions/testscene-run/1.0/internal/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-centos:base
COPY --from=builder /assets /opt/action
