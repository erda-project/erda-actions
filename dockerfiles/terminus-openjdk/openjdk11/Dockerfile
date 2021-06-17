# registry.erda.cloud/erda/openjdk:11
# registry.erda.cloud/erda/openjdk:11.0
# registry.erda.cloud/erda/openjdk:11.0.9
# registry.erda.cloud/erda/openjdk:11.0.9.11
FROM registry.erda.cloud/erda/terminus-centos:7

LABEL maintainer="zhengdong.jzd@alibaba-inc.com"

ENV JAVA_VERSION java11-openjdk
ENV JAVA_HOME /usr/lib/jvm/java
ENV LANG en_US.UTF-8

RUN yum update -y && \
    yum install -y wget && \
    yum install -y java-11-openjdk-11.0.9.11-0.el7_9.x86_64 java-11-openjdk-devel-11.0.9.11-0.el7_9.x86_64 && \
    yum install -y git && \
    yum install -y unzip && \
    yum clean all

# install tools
## greys (~/.greys/lib/)
RUN mkdir /tmp/greys && \
    curl -sf http://ompc.oss.aliyuncs.com/greys/release/greys-stable-bin.zip -o /tmp/greys/greys-stable-bin.zip && \
    pushd /tmp/greys && unzip ./greys-stable-bin.zip && cd greys && sh ./install-local.sh && \
    cp ./greys.sh /bin/greys && \
    popd && rm -rf /tmp/greys
## arthas
RUN mkdir /opt/arthas && \
    curl -sf https://arthas.aliyun.com/arthas-boot.jar -o /opt/arthas/arthas-boot.jar
