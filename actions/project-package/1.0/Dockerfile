FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOPROXY=https://goproxy.cn,direct GOOS=linux GOARCH=$ARCH go build -o /assets/run           actions/project-package/1.0/internal/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /assets /opt/action
