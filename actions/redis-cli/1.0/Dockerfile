FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

# go build
RUN set -x && eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/redis-cli/1.0/internal/cmd"

FROM goodsmileduck/redis-cli

SHELL ["/bin/bash", "--login", "-c"]

COPY --from=builder /opt/action/run /opt/action/run

