FROM registry.erda.cloud/erda-x/golang:1.24 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED=0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"
RUN set -x && eval "go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/java/1.0/internal/cmd"

RUN mkdir -p /opt/action/comp && \
    cp -r actions/java/1.0/comp/* /opt/action/comp

# newest spot agent
RUN bash /opt/action/comp/download_spot_agent.sh
RUN bash /opt/action/comp/download_fonts.sh

FROM registry.erda.cloud/retag/buildkit:v0.11.6 AS buildkit
FROM registry.erda.cloud/retag/docker:20.10.24-cli AS docker-cli
FROM registry.erda.cloud/erda-x/openjdk:8_11_17_21-ubuntu24.04-20260324-114328

ENV HOME=/root

# 设置默认使用 bash 作为 shell
SHELL ["/bin/bash", "-lc"]

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gawk \
        zip \
        unzip \
        curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV MAVEN_VERSION=3.9.12
ENV MAVEN_HOME=/opt/maven
ENV PATH=$MAVEN_HOME/bin:$PATH

RUN curl -fsSL https://erda-project.oss-cn-hangzhou.aliyuncs.com/erda-actions/apache-maven-${MAVEN_VERSION}-bin.tar.gz \
    | tar -xz -C /opt \
 && mv /opt/apache-maven-${MAVEN_VERSION} ${MAVEN_HOME}

# install docker-cli & buildctl
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=docker-cli /usr/local/bin/docker /usr/bin/docker

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
