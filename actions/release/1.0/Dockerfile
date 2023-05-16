FROM registry.erda.cloud/erda-actions/custom-script:2.0 AS builder

ENV CGO_ENABLED 0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN mkdir -p /opt/action/comp/migration/sql && \
    cp -r actions/release/1.0/comp/migration/* /opt/action/comp/migration

RUN go build -o /assets/run github.com/erda-project/erda-actions/actions/release/1.0/internal/cmd

FROM registry.erda.cloud/erda-actions/custom-script:2.0

COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp

#https://github.com/google/bundletool/releases/download/1.11.1/bundletool-all-1.11.1.jar
COPY bundletool-all-1.11.1.jar /opt/action/
