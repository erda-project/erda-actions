#!/usr/bin/env bash
set -eo pipefail
image=registry.erda.cloud/erda/terminus-debian-nu:0.112-node.22.22.lts
# 使用 buildx 构建并同时推送 amd64 和 arm64 双架构镜像
docker buildx build --platform linux/amd64,linux/arm64 -t ${image} --push .
echo "action meta: terminus-debian-nu=$image"
