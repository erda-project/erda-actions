#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.18-n14.19

docker build . -t ${image}
docker push ${image}

echo "action meta: terminus-debian-herd-1.1.18-n14.19=$image"
