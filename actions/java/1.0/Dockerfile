FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.3 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7 AS builder

LABEL maintainer="jiangzhengdong <zhengdong.jzd@alibaba-inc.com>"

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG TARGETARCH

RUN set -x \
    	&& eval "GOPROXY=https://goproxy.cn,direct GOOS=linux GOARCH=$TARGETARCH go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/java/1.0/internal/cmd"

RUN mkdir -p /opt/action/comp && \
    cp -r actions/java/1.0/comp/* /opt/action/comp

# newest spot agent
RUN bash /opt/action/comp/download_spot_agent.sh
RUN bash /opt/action/comp/download_fonts.sh

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/openjdk-maven:8_11-3.9

ARG USER_HOME_DIR="/root"

RUN yum install -y java-11-openjdk-devel wget

RUN yum-config-manager \
         --add-repo \
         http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo \
   && yum install -y docker-ce-cli \
   && yum clean all

ENV MAVEN_HOME /usr/share/maven
ENV MAVEN_CONFIG "$USER_HOME_DIR/.m2"

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
