FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM registry.erda.cloud/retag/golang:1.16-alpine3.14 AS builder

ENV CGO_ENABLED 0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/cmd


FROM registry.erda.cloud/retag/alpine:3.14

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

RUN echo "http://mirrors.aliyun.com/alpine/v3.9/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.9/community/" >> /etc/apk/repositories && \
    apk update && apk add --no-cache docker bash

COPY --from=builder /opt/action/run /opt/action/run
