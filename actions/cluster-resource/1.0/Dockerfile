FROM registry.erda.cloud/erda/terraform:0.12

COPY . /go/src/terminus.io/dice/tools
WORKDIR /go/src/terminus.io/dice/tools
RUN GOOS=linux GOARCH=amd64 go build -o /opt/tools/bin/cluster-resource terminus.io/dice/tools/misc/actions/cluster-resource/1.0/internal/cmd
RUN rm -rf /go/src/terminus.io/

COPY misc/actions/cluster-resource/1.0/internal/run.sh /opt/action/run
COPY scripts/cloudctl/alicloud/ /opt/tools/scripts/cloudctl/alicloud/
COPY bin/cloudctl /opt/tools/bin/cloudctl
COPY templates/ /opt/tools/templates

WORKDIR /opt
