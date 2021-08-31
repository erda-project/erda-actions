#!/usr/bin/env bash

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.8

docker build . -t ${image}
docker push ${image}
