FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/golang:1.17.12-alpine3.15 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run actions/sonar/2.0/internal/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-centos:base AS action

RUN yum install -y wget unzip

RUN mkdir /opt/sonarqube && cd /opt/sonarqube && \
	wget https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/sonarqube/scanner/sonar-scanner-cli-4.4.0.2170-linux.zip && \
    unzip -q sonar-scanner-cli-4.4.0.2170-linux.zip

# nodejs
# https://docs.sonarqube.org/latest/analysis/languages/javascript/
ENV NODE_VERSION 16.16.0
RUN \
    curl --silent --location https://rpm.nodesource.com/setup_16.x | bash - && \
    yum install -y nodejs-$NODE_VERSION
# typescript
RUN npm install -g typescript
ENV NODE_PATH "/usr/lib/node_modules/"

ENV PATH="/opt/sonarqube/sonar-scanner-4.4.0.2170-linux/bin:${PATH}"

COPY --from=builder /assets /opt/action

FROM action
