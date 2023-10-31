#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.26-n18.17

docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push . -f Dockerfile

echo "action meta: terminus-debian-herd-1.1.26-n18.17=$image"
