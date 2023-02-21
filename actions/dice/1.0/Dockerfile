FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run           actions/dice/1.0/internal/run/cmd/*.go
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/when_sigterm  actions/dice/1.0/internal/post/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /assets /opt/action
