FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/gcrane:v0.7.0 as gcrane
FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.3 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/all-in-one:7

ARG TARGETARCH

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=docker/buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx
COPY --from=gcrane /ko-app/gcrane /usr/bin

RUN yum install -y yum-utils device-mapper-persistent-data lvm2 \
  && yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo \
  && sed -i 's+download.docker.com+mirrors.aliyun.com/docker-ce+' /etc/yum.repos.d/docker-ce.repo \
  && yum -y install librdkafka docker-ce-cli

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2 \
  && npm install -g yarn

RUN wget -O nerdctl.tar.gz https://erda-project.oss-cn-hangzhou.aliyuncs.com/erda-actions/nerdctl-1.2.1-linux-${TARGETARCH}.tar.gz \
  && tar zxvf nerdctl.tar.gz -C /usr/bin \
  && rm -f nerdctl.tar.gz

RUN wget https://erda-project.oss-cn-hangzhou.aliyuncs.com/erda-actions/ossutil-v1.7.15-linux-${TARGETARCH} -O /usr/local/bin/ossutil \
  && chmod 755 /usr/local/bin/ossutil

RUN wget -O /usr/local/bin/crictl https://erda-project.oss-cn-hangzhou.aliyuncs.com/k4s/${TARGETARCH}/crictl-v1.24.2 \
  && chmod 755 /usr/local/bin/crictl

WORKDIR $GOPATH
