FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/golang:1.22 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run actions/sonar/2.0/internal/cmd/*.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/golang:1.22 AS action

RUN apt install -y wget unzip

RUN mkdir /opt/sonarqube && cd /opt/sonarqube && \
	wget https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/sonarqube/scanner/sonar-scanner-cli-4.4.0.2170-linux.zip && \
    unzip -q sonar-scanner-cli-4.4.0.2170-linux.zip

# nodejs
# https://docs.sonarqube.org/latest/analysis/languages/javascript/
RUN curl -sL https://deb.nodesource.com/setup_lts.x | bash - && apt update && apt install -y nodejs && rm -fr /var/lib/apt/lists/*
# typescript
RUN npm install -g typescript
ENV NODE_PATH "/usr/lib/node_modules/"

ENV PATH="/opt/sonarqube/sonar-scanner-4.4.0.2170-linux/bin:${PATH}"

COPY --from=builder /assets /opt/action

FROM action
