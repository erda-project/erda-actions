FROM golang:1.15 AS builder

COPY . /go/src/terminus.io/dice/extensions
WORKDIR /go/src/terminus.io/dice/extensions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run actions/dice-version-archive/1.0/internal/cmd/main.go

FROM centos

MAINTAINER chenzhongrun "zhongrun.czr@alibaba-inc.com"

COPY --from=builder /opt/action/run /opt/action/run

RUN chmod 777 /opt/action/*
