FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn go build -o /assets/run actions/release-fetch/1.0/main.go

FROM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /assets /opt/action
