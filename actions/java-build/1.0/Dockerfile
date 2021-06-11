FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

LABEL maintainer="jiangzhengdong <zhengdong.jzd@alibaba-inc.com>"

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN set -x \
    	&& eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/java-build/1.0/internal/cmd"

RUN mkdir -p /opt/action/comp && \
    cp -r actions/java-build/1.0/comp/* /opt/action/comp

# newest spot agent
RUN bash /opt/action/comp/download_spot_agent.sh
RUN bash /opt/action/comp/download_fonts.sh

FROM registry.erda.cloud/erda/terminus-openjdk:v11.0.6

ARG USER_HOME_DIR="/root"

RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
  && curl -fsSL -o /tmp/apache-maven.tar.gz https://mirrors.bfsu.edu.cn/apache/maven/maven-3/3.6.3/binaries/apache-maven-3.6.3-bin.tar.gz \
  && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
  && rm -f /tmp/apache-maven.tar.gz \
  && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

RUN  yum install -y yum-utils java-1.8.0-openjdk-devel-1:1.8.0.272.b10-1.el7_9.x86_64 \
   && yum-config-manager \
         --add-repo \
         http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo \
   && yum install -y docker-ce-cli \
   && yum clean all


ENV MAVEN_HOME /usr/share/maven
ENV MAVEN_CONFIG "$USER_HOME_DIR/.m2"

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
