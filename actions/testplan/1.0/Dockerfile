FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.11.2 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run actions/testplan/1.0/internal/cmd/*.go

FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/pipeline-resource:base AS action
COPY --from=builder /assets /opt/action

FROM action