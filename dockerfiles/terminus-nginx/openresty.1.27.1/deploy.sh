#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/nginx:1.27.1

docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push . -f Dockerfile

echo "action meta: nginx-1.27.1=$image"
