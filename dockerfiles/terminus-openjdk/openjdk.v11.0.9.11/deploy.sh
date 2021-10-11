#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda/openjdk:v11.0.9.11
docker build . -t ${image}
docker push ${image}
echo "action meta: openjdk.v11.0.9.11=${image}"

image2=registry.erda.cloud/erda/openjdk:11
docker tag ${image} ${image2}
docker push ${image2}
echo "action meta: openjdk.11=${image2}"

image3=registry.erda.cloud/erda/openjdk:11.0
docker tag ${image} ${image3}
docker push ${image3}
echo "action meta: openjdk.11.0=${image3}"

image4=registry.erda.cloud/erda/openjdk:11.0.9
docker tag ${image} ${image4}
docker push ${image4}
echo "action meta: openjdk.11.0.9=${image4}"

image5=registry.erda.cloud/erda/openjdk:11.0.9.11
docker tag ${image} ${image5}
docker push ${image5}
echo "action meta: openjdk.11.0.9.11=${image5}"
