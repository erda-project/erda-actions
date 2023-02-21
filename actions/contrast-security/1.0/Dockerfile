FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY .  /erda-actions
WORKDIR /erda-actions

ARG ARCH

# go build
ENV CGO_ENABLED 0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"
RUN export GOPROXY=https://goproxy.io
RUN go mod vendor
RUN set -x && eval "GOOS=linux GOARCH=$ARCH go build $BUILD_FLAGS -o /opt/action/run actions/contrast-security/1.0/cmd/main.go"

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /opt/action/run /opt/action/run