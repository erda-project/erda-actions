#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-alpine-herd:1.1.11-n14.17

docker build . -t ${image}
docker push ${image}

echo "action meta: terminus-alpine-herd-1.1.11-n14.17=$image"
