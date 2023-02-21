FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH go build -o /opt/action/run github.com/erda-project/erda-actions/actions/lib-publish/1.0/internal/cmd

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /opt/action/run /opt/action/run
