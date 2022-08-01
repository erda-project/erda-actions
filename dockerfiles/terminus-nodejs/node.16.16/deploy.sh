#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-node:16.16-lts

docker buildx build --platform linux/amd64 -t ${image} --push . -f Dockerfile.debian.npm.8.11

echo "action meta: terminus-debian-node=$image"
