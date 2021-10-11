#!/usr/bin/env bash
set -eo pipefail

image=registry.erda.cloud/erda/terminus-tomcat:v7.0.96-jdk8
docker build . -t ${image}
docker push ${image}
echo "action meta: terminus-tomcat.v7.0.96-jdk8=${image}"
