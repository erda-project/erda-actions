FROM registry.erda.cloud/erda/terminus-golang:1.16.7 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOPROXY=https://goproxy.io,direct GOOS=linux GOARCH=amd64 go build -o /assets/run actions/semgrep/1.0/internal/cmd/*.go

FROM registry.erda.cloud/erda/terminus-centos:base AS action

RUN yum install -y python3
RUN python3 -m pip install semgrep  -i https://mirrors.aliyun.com/pypi/simple

COPY --from=builder /assets /opt/action

FROM action
