FROM registry.erda.cloud/retag/openjdk:8u181

ENV JAVA_VERSION openjdk-8u181
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8

# timezone & locate
ENV TZ "Asia/Shanghai"
RUN apt update && apt install tzdata locales -y \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "zh_CN.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen

# install tools
## greys
RUN curl -sLk https://ompc.oss.aliyuncs.com/greys/install.sh | bash && cp ./greys.sh /bin/greys && (greys || true)

## arthas
RUN mkdir /opt/arthas && \
    curl -sf https://arthas.aliyun.com/arthas-boot.jar -o /opt/arthas/arthas-boot.jar
