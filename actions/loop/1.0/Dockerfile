FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.17.12 AS builder

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

RUN set -x \
    	&& eval "GOOS=linux GOARCH=$ARCH go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/loop/1.0/internal"

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base

COPY --from=builder /opt/action/run /opt/action/run

CMD /opt/action/run
