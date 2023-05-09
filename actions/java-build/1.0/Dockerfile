FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/java-build/1.0/internal/cmd

RUN mkdir -p /opt/action/comp && \
    cp -r actions/java-build/1.0/comp/* /opt/action/comp

# newest spot agent
RUN bash /opt/action/comp/download_spot_agent.sh
RUN bash /opt/action/comp/download_fonts.sh

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0

RUN alternatives --set java $(alternatives --list | grep java_sdk_11  | awk '{print $3}' | head -n 1)/bin/java
RUN alternatives --set javac $(alternatives --list | grep java_sdk_11  | awk '{print $3}' | head -n 1)/bin/javac
ENV JAVA_HOME=/usr/lib/jvm/java-11

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
