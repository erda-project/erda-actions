FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM registry.erda.cloud/erda/git-golang-image:1.1

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV GOPROXY=https://goproxy.cn/
ENV GO111MODULE=on

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0

RUN apk add --no-cache librdkafka-dev bash libc6-compat  \
    openjdk8 nodejs npm yarn docker-cli

RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
  && curl -fsSL -o /tmp/apache-maven.tar.gz https://mirrors.bfsu.edu.cn/apache/maven/maven-3/3.6.3/binaries/apache-maven-3.6.3-bin.tar.gz \
  && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
  && rm -f /tmp/apache-maven.tar.gz \
  && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

ADD actions/custom-script/1.0/assets/settings.xml /root/.m2/settings.xml

RUN ln -sf /bin/bash /bin/sh

RUN wget http://gosspublic.alicdn.com/ossutil/1.5.0/ossutil64 -O /usr/local/bin/ossutil && chmod 755 /usr/local/bin/ossutil

WORKDIR $GOPATH



