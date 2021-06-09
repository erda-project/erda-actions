FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.14 as builder

MAINTAINER shenli shenli@terminus.io

ENV GOLANG_VERSION 1.12.7

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN mkdir -p /opt/action/comp && \
    cp -r actions/golang/1.0/comp/* /opt/action/comp

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/golang/1.0/internal/cmd/main.go

FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-golang:1.14
RUN yum install -y docker make
RUN yum install -y gcc
RUN yum install -y g++
RUN yum install -y gcc-c++
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp

