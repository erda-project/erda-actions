#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda/terminus-nodejs:12.13

docker build . -t ${image} -f Dockerfile.node.12.13
docker push ${image}

echo "action meta: terminus-nodejs-12.13=$image"
