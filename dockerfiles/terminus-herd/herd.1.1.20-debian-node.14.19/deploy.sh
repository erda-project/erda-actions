#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.20-n14.19

docker buildx build --platform linux/amd64 -t ${image} --push . -f Dockerfile

echo "action meta: terminus-debian-herd-1.1.20-n14.19=$image"
