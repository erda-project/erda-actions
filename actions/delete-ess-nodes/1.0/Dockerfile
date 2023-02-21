FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

RUN GOOS=linux GOARCH=$ARCH go build -o /opt/action/run github.com/erda-project/erda-actions/actions/delete-ess-nodes/1.0/internal

WORKDIR /opt
