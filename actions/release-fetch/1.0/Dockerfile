FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH GOPROXY=https://goproxy.cn go build -o /assets/run actions/release-fetch/1.0/main.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /assets /opt/action
