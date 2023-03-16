FROM --platform=$BUILDPLATFORM registry.erda.cloud/erda-x/gcrane:v0.7.0 as gcrane
FROM --platform=$BUILDPLATFORM registry.erda.cloud/erda-x/buildkit:v0.11.3 as buildkit
FROM --platform=$BUILDPLATFORM registry.erda.cloud/erda-x/terminus-golang:1.19.7

ARG TARGETARCH

# install packages
RUN apk add --update git openssh make curl libgit2-dev pkgconfig gcc musl-dev sed python3 build-base binutils-gold

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV GOPROXY=https://goproxy.cn/
ENV GO111MODULE=on

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2

RUN apk add --no-cache librdkafka-dev bash libc6-compat  \
    openjdk8 nodejs npm yarn docker-cli tzdata

RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
  && curl -fsSL -o /tmp/apache-maven.tar.gz https://mirrors.bfsu.edu.cn/apache/maven/maven-3/3.9.0/binaries/apache-maven-3.9.0-bin.tar.gz  \
  && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
  && rm -f /tmp/apache-maven.tar.gz \
  && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

ADD actions/custom-script/1.0/assets/settings.xml /root/.m2/settings.xml

RUN ln -sf /bin/bash /bin/sh

RUN wget -O nerdctl-1.2.1-linux-${TARGETARCH}.tar.gz https://github.com/containerd/nerdctl/releases/download/v1.2.1/nerdctl-1.2.1-linux-${TARGETARCH}.tar.gz \
    && tar zxvf nerdctl-1.2.1-linux-${TARGETARCH}.tar.gz -C /usr/bin \
    && rm -f nerdctl-1.2.1-linux-${TARGETARCH}.tar.gz

RUN #wget http://gosspublic.alicdn.com/ossutil/1.7.8/ossutil64 -O /usr/local/bin/ossutil && chmod 755 /usr/local/bin/ossutil

RUN wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.24.1/crictl-v1.24.1-linux-${TARGETARCH}.tar.gz \
    && tar zxvf crictl-v1.24.1-linux-${TARGETARCH}.tar.gz -C /usr/local/bin \
    && rm -f crictl-v1.24.1-linux-${TARGETARCH}.tar.gz

COPY --from=docker/buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx

COPY --from=gcrane /ko-app/gcrane /usr/bin

WORKDIR $GOPATH



