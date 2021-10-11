#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda/terminus-openjdk:v1.8.0.272
docker build . -t ${image}
docker push ${image}
echo "action meta: terminus-openjdk.v1.8.0.272=${image}"

image2=registry.erda.cloud/erda/openjdk:8
docker tag ${image} ${image2}
docker push ${image2}
echo "action meta: openjdk.8=${image2}"

image3=registry.erda.cloud/erda/openjdk:8u272
docker tag ${image} ${image3}
docker push ${image3}
echo "action meta: openjdk.272=${image3}"
