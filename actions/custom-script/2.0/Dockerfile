FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/gcrane:v0.7.0 as gcrane
FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.6 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/all-in-one:7

ARG TARGETARCH

COPY --from=gcrane /ko-app/gcrane /usr/bin
COPY --from=docker/buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

RUN yum install -y yum-utils device-mapper-persistent-data lvm2 \
  && yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo \
  && sed -i 's+download.docker.com+mirrors.aliyun.com/docker-ce+' /etc/yum.repos.d/docker-ce.repo \
  && yum -y install librdkafka docker-ce-cli-19.03.15 \
  && yum clean all

RUN wget -O nerdctl.tar.gz https://erda-project.oss-cn-hangzhou.aliyuncs.com/erda-actions/nerdctl-1.2.1-linux-${TARGETARCH}.tar.gz \
  && tar -xzvf nerdctl.tar.gz -C /usr/bin \
  && rm -fv nerdctl.tar.gz

RUN wget -O /usr/local/bin/ossutil https://erda-project.oss-cn-hangzhou.aliyuncs.com/erda-actions/ossutil-v1.7.15-linux-${TARGETARCH} \
  && chmod 755 /usr/local/bin/ossutil

RUN wget -O /usr/local/bin/crictl https://erda-project.oss-cn-hangzhou.aliyuncs.com/k4s/${TARGETARCH}/crictl-v1.24.2 \
  && chmod 755 /usr/local/bin/crictl

COPY actions/custom-script/2.0/assets/settings.xml /root/.m2/settings.xml
COPY actions/custom-script/2.0/assets/settings.xml /home/dice/.m2/settings.xml

RUN chown -R dice:dice /home/dice/.m2 \
  && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2 \
  && npm install -g yarn

ARG TARGETARCH
RUN yum erase -y git && cd /opt && curl -O -L https://erda-project.oss-cn-hangzhou.aliyuncs.com/erda-addons/git-$TARGETARCH.tgz && tar -xzvf git-$TARGETARCH.tgz && rm -fv git-$TARGETARCH.tgz
ENV PATH=/opt/git/bin:$PATH

WORKDIR $GOPATH
