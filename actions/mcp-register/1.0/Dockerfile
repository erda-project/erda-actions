FROM registry.erda.cloud/erda-x/golang:1.24 as builder

ENV CGO_ENABLED=0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /assets/run           actions/mcp-register/1.0/internal/cmd/*.go

FROM registry.erda.cloud/erda-x/oraclelinux:7
COPY --from=builder /assets /opt/action
