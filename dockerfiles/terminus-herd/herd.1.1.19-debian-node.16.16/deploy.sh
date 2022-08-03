#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.19-n16.16

docker buildx build --platform linux/amd64 -t ${image} --push . -f Dockerfile

echo "action meta: terminus-debian-herd-1.1.19-n16.16=$image"
