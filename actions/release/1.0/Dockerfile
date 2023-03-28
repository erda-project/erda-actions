FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.3 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7 AS builder

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"
ENV GOPROXY=https://goproxy.cn

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG TARGETARCH

RUN mkdir -p /opt/action/comp/migration/sql && \
    cp -r actions/release/1.0/comp/migration/* /opt/action/comp/migration

RUN set -x \
    	&& eval "GOOS=linux GOARCH=$TARGETARCH go build $BUILD_FLAGS -o /assets/run github.com/erda-project/erda-actions/actions/release/1.0/internal/cmd"

FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/alpine:3.14

RUN apk update && apk add --no-cache docker openjdk8-jre
RUN wget https://github.com/google/bundletool/releases/download/1.11.1/bundletool-all-1.11.1.jar && mkdir /opt/action \
    && cp ./bundletool-all-1.11.1.jar /opt/action/
RUN apk add --no-cache bash tar libarchive-tools
RUN cp  /usr/bin/bsdtar /usr/bin/tar

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
