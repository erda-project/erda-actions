FROM registry.erda.cloud/retag/buildkit:v0.11.6 AS buildkit

FROM registry.erda.cloud/erda-actions/custom-script:2.0 AS custom-script

FROM registry.erda.cloud/erda-x/golang:1.22 AS builder

COPY . /go/src/github.com/erda-project/erda-actions

WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js/1.0/internal/cmd && \
    mkdir -p /opt/action/comp && \
    cp -r actions/js/1.0/comp/* /opt/action/comp

FROM registry.erda.cloud/erda-x/nvm:0.40

ENV NODE_OPTIONS=--max_old_space_size=1800
# install buildkit
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
COPY --from=custom-script /usr/bin/docker /usr/bin/docker


