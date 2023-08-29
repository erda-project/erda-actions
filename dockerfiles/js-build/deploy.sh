#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-js-build:2.0

docker buildx build --platform linux/amd64 -t ${image} --push . -f Dockerfile

echo "action meta: terminus-js-build=$image"
