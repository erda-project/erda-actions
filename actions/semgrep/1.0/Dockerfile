FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script-action:2.0-20230608153022-685950d AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN GOPROXY=https://goproxy.io,direct go build -o /assert/run actions/semgrep/1.0/internal/cmd/*.go

FROM python:3.11-bullseye as action

RUN python3 -m pip install --upgrade semgrep  -i https://mirrors.aliyun.com/pypi/simple

COPY --from=builder /assert /opt/action
