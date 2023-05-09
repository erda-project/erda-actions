FROM registry.erda.cloud/erda-x/golang:1 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions/actions
RUN rm -rf `ls | grep -v erda-mysql-migration`
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go mod tidy -compat=1.17

RUN go build -o /opt/action/run actions/erda-mysql-migration/1.0-57/internal/action/cmd/main.go
RUN go build -o /opt/action/erda-migrate actions/erda-mysql-migration/1.0-57/internal/local/cmd/main.go

FROM registry.erda.cloud/erda-actions/erda-mysql-migration-sandbox:80

VOLUME ["/log", "/migrations"]
ENV MIGRATION_DIR=/migrations

COPY --from=builder /go/src/github.com/erda-project/erda-actions/go.mod /opt/action/go.mod
COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/erda-migrate /usr/bin/erda-migrate

RUN chmod 755 /opt/action/run
RUN chmod 755 /usr/bin/erda-migrate

CMD ["erda-migrate"]
