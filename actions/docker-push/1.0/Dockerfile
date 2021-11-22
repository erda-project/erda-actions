FROM registry.erda.cloud/retag/golang:1.16-alpine3.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/docker-push/1.0/internal

FROM registry.erda.cloud/retag/gcrane:v0.7.0 as gcrane
FROM registry.erda.cloud/retag/alpine:3.14

COPY --from=gcrane /ko-app/gcrane /usr/bin

RUN echo "http://mirrors.aliyun.com/alpine/v3.9/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.9/community/" >> /etc/apk/repositories && \
    apk update && apk add --no-cache docker bash

COPY --from=builder /opt/action/run /opt/action/run
