#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-node:12.22-lts

docker build . -t ${image} -f Dockerfile.debian.npm.6.14
docker push ${image}

echo "action meta: terminus-debian-node-12.22=$image"
