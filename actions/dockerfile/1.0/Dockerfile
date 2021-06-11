FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

ENV CGO_ENABLED 0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/cmd


FROM registry.erda.cloud/erda/terminus-alpine:base

RUN echo "http://mirrors.aliyun.com/alpine/v3.9/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.9/community/" >> /etc/apk/repositories && \
    apk update && apk add --no-cache docker

COPY --from=builder /opt/action/run /opt/action/run
