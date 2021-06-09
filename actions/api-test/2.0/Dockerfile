FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run actions/api-test/2.0/internal/*.go

FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-centos:base
RUN yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm && yum install -y jq
COPY --from=builder /assets /opt/action
