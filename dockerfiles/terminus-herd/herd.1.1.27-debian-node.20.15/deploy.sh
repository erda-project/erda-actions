#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.27-n20.15

docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push . -f Dockerfile

echo "action meta: terminus-debian-herd-1.1.27-n20.15=$image"
