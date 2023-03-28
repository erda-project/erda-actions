FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.3 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7 as builder

MAINTAINER shenli shenli@terminus.io

ARG TARGETARCH

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN mkdir -p /opt/action/comp && \
    cp -r actions/golang/1.0/comp/* /opt/action/comp

# go build
RUN GOPROXY=https://goproxy.cn,direct GOOS=linux GOARCH=$TARGETARCH go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/golang/1.0/internal/cmd/main.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7
RUN yum install -y docker make gcc g++ gcc-c++
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

