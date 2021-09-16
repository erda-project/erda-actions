#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.10-n12.22

docker build . -t ${image}
docker push ${image}

echo "action meta: terminus-debian-herd-1.1.10-n12.22=$image"
