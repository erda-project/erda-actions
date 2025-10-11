FROM registry.erda.cloud/erda-x/golang:1.24 AS builder

WORKDIR /go/src/github.com/erda-project/erda-actions

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/cmd

FROM registry.erda.cloud/retag/buildkit:v0.11.6 AS buildkit
FROM registry.erda.cloud/retag/docker:20.10.24-cli AS docker-cli
FROM registry.erda.cloud/erda-x/debian-bookworm:12

# install docker-cli & buildctl
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=docker-cli /usr/local/bin/docker /usr/bin/docker

COPY --from=builder /opt/action/run /opt/action/run
RUN chmod +x /opt/action/run
