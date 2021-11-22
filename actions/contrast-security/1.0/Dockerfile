FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY .  /erda-actions
WORKDIR /erda-actions

# go build
ENV CGO_ENABLED 0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"
RUN export GOPROXY=https://goproxy.io
RUN go mod vendor
RUN set -x && eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run actions/contrast-security/1.0/cmd/main.go"

FROM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /opt/action/run /opt/action/run