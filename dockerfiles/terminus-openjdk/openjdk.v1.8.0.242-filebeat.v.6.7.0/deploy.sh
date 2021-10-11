#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda/terminus-openjdk:v1.8.0.242-filebeat.v6.7.0

docker build . -t ${image}
docker push ${image}

echo "action meta: terminus-openjdk-v1.8.0.242-filebeat.v6.7.0=${image}"
