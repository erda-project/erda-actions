FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/lib-publish/1.0/internal/cmd

FROM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /opt/action/run /opt/action/run
