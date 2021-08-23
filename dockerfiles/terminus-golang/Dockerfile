# output image: registry.cn-hangzhou.aliyuncs.com/terminus/terminus-golang:1.14

FROM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-centos:base

MAINTAINER linjun lj@terminus.io

ARG GOLANG_VERSION
ARG GO_REL_SHA256

RUN set -eux; \
    \
    goRelArch='linux-amd64';\
    \
    url="https://golang.google.cn/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
    curl -k -L -o go.tgz "$url"; \
    echo "${GO_REL_SHA256} *go.tgz" | sha256sum -c -; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    \
    export PATH="/usr/local/go/bin:$PATH"; \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH