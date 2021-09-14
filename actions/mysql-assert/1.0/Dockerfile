FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY .  /erda-actions
WORKDIR /erda-actions

ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

RUN export GOPROXY=https://goproxy.io
RUN go mod vendor
# go build
RUN set -x && eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run actions/mysql-assert/1.0/internal/cmd/main.go"

FROM oraclelinux:7-slim

ARG KEY=https://repo.mysql.com/RPM-GPG-KEY-mysql
ARG REPO=https://repo.mysql.com

ARG PACKAGE_URL_SHELL=$REPO/yum/mysql-tools-community/el/7/x86_64/mysql-shell-8.0.13-1.el7.x86_64.rpm

# Install server
RUN rpmkeys --import $KEY \
  && yum install -y $PACKAGE_URL_SHELL libpwquality \
  && yum clean all \
  && mkdir /docker-entrypoint-initdb.d
RUN yum install -y jq

COPY --from=builder /opt/action/run /opt/action/run
