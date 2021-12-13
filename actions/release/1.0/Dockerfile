FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM registry.erda.cloud/retag/golang:1.16-alpine3.14 AS builder

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN mkdir -p /opt/action/comp/migration/sql && \
    cp -r actions/release/1.0/comp/migration/* /opt/action/comp/migration

RUN set -x \
    	&& eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /assets/run github.com/erda-project/erda-actions/actions/release/1.0/internal/cmd"

FROM registry.erda.cloud/retag/alpine:3.14

RUN echo "http://mirrors.aliyun.com/alpine/v3.12/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.12/community/" >> /etc/apk/repositories && \
    apk update && apk add --no-cache docker
RUN apk add --no-cache bash tar

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
