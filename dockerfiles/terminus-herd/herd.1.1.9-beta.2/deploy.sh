#!/usr/bin/env bash

image=erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/terminus-herd:1.1.9-beta.2

docker build . -t ${image}
docker push ${image}
