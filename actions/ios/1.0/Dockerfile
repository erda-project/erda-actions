FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.14 as builder

MAINTAINER shenli shenli@terminus.io

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/ios/1.0/internal/cmd/main.go

FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/pipeline-resource:base
COPY --from=builder /assets /opt/action