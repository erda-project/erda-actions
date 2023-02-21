FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY .  /erda-actions
WORKDIR /erda-actions

ARG ARCH

ENV CGO_ENABLED 0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"
RUN export GOPROXY=https://goproxy.io
RUN go mod vendor
# go build
RUN set -x && eval "GOOS=linux GOARCH=$ARCH go build $BUILD_FLAGS -o /opt/action/run actions/jsonparse/1.0/internal/cmd/main.go"

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-alpine:base

RUN apk add jq
RUN apk add --update nodejs nodejs-npm
RUN npm i -g jackson-converter@1.0.10

COPY --from=builder /opt/action/run /opt/action/run
