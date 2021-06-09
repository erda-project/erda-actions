FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run actions/integration-test/1.0/internal/cmd/*.go

FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/pipeline-resource:base AS action
COPY actions/integration-test/1.0/internal/settings.xml /root/.m2/settings.xml
COPY --from=builder /assets /opt/action
