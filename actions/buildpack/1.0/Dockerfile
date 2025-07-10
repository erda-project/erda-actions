FROM registry.erda.cloud/retag/buildkit:v0.11.3 AS buildkit
FROM registry.erda.cloud/retag/docker:28.2.2-cli AS docker-cli
FROM registry.erda.cloud/erda-x/golang:1.22 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /assets/run github.com/erda-project/erda-actions/actions/buildpack/1.0/internal

WORKDIR actions/buildpack/1.0
RUN bash sync_assets.sh

FROM registry.erda.cloud/erda-x/debian-bookworm:12

RUN apt-get update && \
    apt-get install --no-install-recommends -y procps && \
    rm -rf /var/lib/apt/lists/*

COPY --from=docker-cli /usr/local/bin/docker /usr/bin/docker
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

COPY --from=builder /assets /opt/action
COPY actions/buildpack/1.0/bp /opt/action/bp
