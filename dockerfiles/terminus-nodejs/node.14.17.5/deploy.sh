#!/usr/bin/env bash

image=registry.erda.cloud/erda-actions/terminus-debian-node:14.17-lts

docker build . -t ${image}
docker push ${image}

echo "action meta: terminus-debian-node=$image"