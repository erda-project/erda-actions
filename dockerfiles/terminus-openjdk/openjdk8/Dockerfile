# registry.erda.cloud/erda/openjdk:8
# registry.erda.cloud/erda/openjdk:8u272
FROM registry.erda.cloud/erda/terminus-centos:7

LABEL maintainer="xgs@terminus.io"

ENV JAVA_VERSION openjdk-8u272
ENV JAVA_HOME /usr/lib/jvm/java
ENV LANG en_US.UTF-8

RUN yum update -y && \
    yum install -y wget && \
    yum install -y java-1.8.0-openjdk-1.8.0.272.b10-1.el7_9.x86_64 java-1.8.0-openjdk-devel-1.8.0.272.b10-1.el7_9.x86_64 && \
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
