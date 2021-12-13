FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/gitbook/1.0/internal/cmd
RUN mkdir -p /opt/action/comp && \
    cp -r actions/gitbook/1.0/comp/* /opt/action/comp

FROM registry.erda.cloud/erda/terminus-nodejs:12.13

RUN curl -sSL -q -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo && \
    curl -sSL -q -o /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo && \
    sed -i -e '/mirrors.cloud.aliyuncs.com/d' -e '/mirrors.aliyuncs.com/d' /etc/yum.repos.d/CentOS-Base.repo && \
    yum clean all && yum makecache && yum install -y docker && npm install -g gitbook-cli

RUN gitbook install highlight search sharing font-settings && rm -rf /tmp/*
#RUN cd /tmp && gitbook build && rm -rf _book

#highlight search sharing font-settings
COPY --from=builder /opt/action/run /opt/action/run
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /opt/action/comp /opt/action/comp
