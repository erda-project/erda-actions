FROM registry.erda.cloud/erda/terminus-golang:1.13 AS builder

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN set -x \
    	&& eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal"

RUN mkdir -p /opt/action/assets && \
    cp -r actions/java-deploy/1.0/internal/assets/* /opt/action/assets

FROM registry.erda.cloud/erda/terminus-maven:3-jdk-8-alpine

RUN echo "http://mirrors.aliyun.com/alpine/v3.9/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.9/community/" >> /etc/apk/repositories
#    apk update && apk add --no-cache docker

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/assets /opt/action/assets
