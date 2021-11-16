FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js/1.0/internal/cmd
RUN mkdir -p /opt/action/comp && \
    cp -r actions/js/1.0/comp/* /opt/action/comp

FROM registry.erda.cloud/erda/terminus-nodejs:12.13

RUN yum install -y yum-utils && \
    yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo && \
    yum -y install docker-ce-cli

ENV NODE_OPTIONS=--max_old_space_size=1800

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /opt/action/comp /opt/action/comp
