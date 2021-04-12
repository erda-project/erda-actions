FROM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

# go build
RUN set -x && eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/mysql-cli/1.0/internal/cmd"

FROM imega/mysql-client

COPY --from=builder /opt/action/run /opt/action/run