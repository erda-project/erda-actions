# Author: hustcer
# Created at 2023-08-29T11:08:18+08:00
# Run the following build command from the root directory of this repository:
# docker build . -t terminus/js-build-2.0 -f dockerfiles/js-build/Dockerfile

FROM --platform=$TARGETPLATFORM golang:1.21.0-bookworm AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js-build/1.0/internal/cmd

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/terminus-debian-node:18.17-lts

LABEL maintainer=hustcer<majun@terminus.io>

COPY --from=builder /opt/action/run /opt/action/run

RUN chmod 755 /opt/action/*

SHELL ["/bin/sh", "--login", "-c"]

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash

COPY actions/js-build/1.0/install.sh /tmp/install.sh

RUN bash /tmp/install.sh && rm -fv /tmp/install.sh \
    && source ~/.bashrc && nvm install 14 && nvm install 16 && nvm install 18 && nvm use 18 \
    && npm i vite rollup webpack webpack-cli -g && npm i -g cnpm --registry=https://registry.npmmirror.com
