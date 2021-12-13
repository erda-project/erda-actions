ARG GO_VERSION
FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit

FROM registry.erda.cloud/erda/terminus-golang:${GO_VERSION} as builder

MAINTAINER shenli shenli@terminus.io

ENV GOLANG_VERSION ${GO_VERSION}

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN mkdir -p /opt/action/comp && \
    cp -r actions/golang/1.0/comp/* /opt/action/comp

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/golang/1.0/internal/cmd/main.go

FROM registry.erda.cloud/erda/terminus-golang:${GO_VERSION}
RUN yum install -y docker make gcc g++ gcc-c++
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

