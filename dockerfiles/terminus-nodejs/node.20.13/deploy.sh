#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-node:20.13-lts

docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push . -f Dockerfile.debian.npm.10.5

echo "action meta: terminus-debian-node=$image"
