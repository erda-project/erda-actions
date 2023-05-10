FROM --platform=$BUILDPLATFORM registry.erda.cloud/retag/gcrane:v0.15.1 as gcrane
FROM --platform=$BUILDPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7 as builder

ARG TARGETARCH

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

ENV GOPROXY=https://goproxy.cn,direct

RUN GOOS=linux GOARCH=$TARGETARCH go build -o /opt/action/run github.com/erda-project/erda-actions/actions/docker-push/1.0/internal

FROM --platform=$BUILDPLATFORM registry.erda.cloud/retag/debian:buster-slim

COPY --from=gcrane /ko-app/gcrane /usr/bin
COPY --from=builder /opt/action/run /opt/action/run
