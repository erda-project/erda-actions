FROM registry.erda.cloud/erda/terminus-golang:1.14

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOOS=linux GOARCH=amd64 go build -o /opt/action/run github.com/erda-project/erda-actions/actions/ess-info/1.0/internal/cmd

WORKDIR /opt
