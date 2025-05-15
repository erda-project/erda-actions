FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/golang:1.22 AS builder

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

FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.6 AS buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/openjdk:8_11_17_21

ENV HOME=/root

# 设置默认使用 bash 作为 shell
SHELL ["/bin/bash", "-lc"]

# install maven by sdkman
RUN dnf install -y zip && dnf clean all
RUN curl -s "https://get.sdkman.io" | bash
RUN source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install maven 3.9.9
ENV PATH="$HOME/.sdkman/candidates/maven/current/bin:$PATH"

# install docker-cli & buildctl
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
RUN dnf install -y dnf-plugins-core \
    && dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo \
    && dnf install -y docker-ce-cli \
    && dnf clean all

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
