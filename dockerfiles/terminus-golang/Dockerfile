# output image: registry.cn-hangzhou.aliyuncs.com/terminus/terminus-golang:1.14

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/oraclelinux:7

MAINTAINER linjun lj@terminus.io

ARG TARGETARCH

ARG GOLANG_VERSION

RUN set -eux; \
    goRelArch='linux-amd64';\
    url="https://golang.google.cn/dl/go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz"; \
    curl -k -L -o go.tgz "$url"; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    export PATH="/usr/local/go/bin:$PATH"; \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH