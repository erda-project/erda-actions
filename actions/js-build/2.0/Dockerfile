# Author: hustcer
# Created at 2023-08-29T11:08:18+08:00
# Tools from terminus-debian-node:18.17-lts:
#   Git version 2.39.2, gcc: 12.2.0, curl: 7.88.1, Wget: 1.21.3
#   Node: 18.17.1 LTS, npm: 9.6.7, yarn: 1.22.19, pnpm: 8.6.12, fnm: 1.35.1
# Newly Added:
#   nvm: 0.39.5, vite: 4.4.9, rollup: 3.28.1, webpack: 5.88.2, webpack-cli: 5.1.4
#   Node: v14.21.3 / v16.20.2 / v18.17.1, and default v18.17.1, cnpm: 9.2.0
# Run the following build command from the root directory of this repository:
# docker build . -t terminus/js-build-2.0 -f actions/js-build/1.0/Dockerfile-2.0

FROM --platform=$TARGETPLATFORM golang:1.21.0-bookworm AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js-build/1.0/internal/cmd

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/terminus-debian-node:18.17-lts

LABEL maintainer=hustcer<majun@terminus.io>

COPY --from=builder /opt/action/run /opt/action/run

RUN chmod 755 /opt/action/*

SHELL ["/bin/sh", "--login", "-c"]

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

COPY actions/js-build/1.0/install.sh /tmp/install.sh

RUN bash /tmp/install.sh && rm -fv /tmp/install.sh \
    && source ~/.bashrc && nvm install 14 && nvm install 16 && nvm install 18 && nvm install 20 && nvm use 20 \
    && npm i vite rollup webpack webpack-cli -g && npm i -g cnpm --registry=https://registry.npmmirror.com
