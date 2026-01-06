FROM registry.erda.cloud/erda-x/golang:1.24 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED=0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"
RUN set -x && eval "go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/java-build/1.0/internal/cmd"

RUN mkdir -p /opt/action/comp && \
    cp -r actions/java-build/1.0/comp/* /opt/action/comp

RUN bash /opt/action/comp/download_spot_agent.sh
RUN bash /opt/action/comp/download_fonts.sh

FROM registry.erda.cloud/erda-x/openjdk:8_11_17_21-ubuntu24.04-v2

ENV MAVEN_VERSION=3.9.12
ENV MAVEN_HOME=/opt/maven
ENV PATH=$MAVEN_HOME/bin:$PATH

RUN curl -fsSL https://archive.apache.org/dist/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz \
    | tar -xz -C /opt \
 && mv /opt/apache-maven-${MAVEN_VERSION} ${MAVEN_HOME}

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
