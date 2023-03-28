FROM --platform=$BUILDPLATFORM registry.erda.cloud/retag/buildkit:v0.11.3 as buildkit
FROM --platform=$BUILDPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7 AS builder

ENV CGO_ENABLED 0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG TARGETARCH

RUN GOOS=linux GOARCH=$TARGETARCH go build -o /opt/action/run github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/cmd

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

RUN yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo \
    && yum install -y docker-ce-cli
