#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda/terminus-openjdk:v1.8.0.181-filebeat.v5.6.16

docker build . -t ${image}
docker push ${image}

echo "action meta: terminus-openjdk-v1.8.0.181-filebeat.v5.6.16=${image}"
