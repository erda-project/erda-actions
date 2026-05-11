#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-node:22.22-lts

docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push . -f Dockerfile.debian.npm.10.9

echo "action meta: terminus-debian-node=$image"
