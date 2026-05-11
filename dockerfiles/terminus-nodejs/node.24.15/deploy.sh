#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-node:24.15-lts

docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push . -f Dockerfile.debian.npm.11.12

echo "action meta: terminus-debian-node=$image"
