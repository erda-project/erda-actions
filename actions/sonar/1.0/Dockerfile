FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOOS=linux GOARCH=amd64 go build -o /assets/run actions/sonar/1.0/internal/*.go

FROM registry.erda.cloud/erda/terminus-centos:base AS action

RUN yum install -y wget unzip

RUN mkdir /opt/sonarqube && cd /opt/sonarqube && \
	wget https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/sonarqube/scanner/sonar-scanner-cli-4.4.0.2170-linux.zip && \
    unzip -q sonar-scanner-cli-4.4.0.2170-linux.zip

# nodejs
# https://docs.sonarqube.org/latest/analysis/languages/javascript/
ENV NODE_VERSION 12.13.1
RUN \
    curl --silent --location https://rpm.nodesource.com/setup_12.x | bash - && \
    yum install -y nodejs-$NODE_VERSION
# typescript
RUN npm install -g typescript
ENV NODE_PATH "/usr/lib/node_modules/"

ENV PATH="/opt/sonarqube/sonar-scanner-4.4.0.2170-linux/bin:${PATH}"

COPY --from=builder /assets /opt/action

FROM action
