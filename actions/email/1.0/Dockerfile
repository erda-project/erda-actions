FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN set -x \
    	&& eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/email/1.0/internal/cmd"

FROM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-golang:1.14

COPY --from=builder /opt/action/run /opt/action/run
