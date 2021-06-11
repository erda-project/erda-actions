FROM registry.erda.cloud/erda/terminus-golang:1.11.2 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run actions/unit-test/1.0/internal/cmd/*.go

FROM registry.erda.cloud/erda/pipeline-resource:base AS action
COPY actions/unit-test/1.0/internal/settings.xml /root/.m2/settings.xml
COPY actions/unit-test/1.0/internal/init.gradle /root/.gradle/init.gradle
COPY --from=builder /assets /opt/action

RUN GOOS=linux GOARCH=amd64 go get github.com/jstemmer/go-junit-report

FROM action
