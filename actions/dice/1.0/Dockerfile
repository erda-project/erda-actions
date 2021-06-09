FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run           actions/dice/1.0/internal/run/cmd/*.go
RUN GOOS=linux GOARCH=amd64 go build -o /assets/when_sigterm  actions/dice/1.0/internal/post/cmd/*.go

FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-centos:base
COPY --from=builder /assets /opt/action
