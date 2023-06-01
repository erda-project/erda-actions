FROM registry.erda.cloud/erda-actions/custom-script:2.0 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/js/1.0/internal/cmd
RUN mkdir -p /opt/action/comp && \
    cp -r actions/js/1.0/comp/* /opt/action/comp

FROM registry.erda.cloud/erda-actions/custom-script:2.0

ENV NODE_OPTIONS=--max_old_space_size=1800

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
