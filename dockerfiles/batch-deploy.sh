#!/bin/bash
set -e
set -u
set -o pipefail

# refer to erda.cloud erda-release repo, pipeline: deploy-action-dockerfiles.yml

registry="${DOCKER_REGISTRY-}"
username="${DOCKER_USERNAME-}"
password="${DOCKER_PASSWORD-}"

if [[ -z $registry ]]; then
    echo No registry specified! You should specify it through env "'DOCKER_REGISTRY'".
    exit 1
fi

if [[ -n $username && -n $password ]]; then
    echo "You specifed docker username&password, execute login step."
    docker login ${DOCKER_REGISTRY} -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}
fi

# find all deploy.sh and deploy
deployFiles=$(find "$(cd $(dirname "$0"); pwd)" -type f -name deploy.sh)
for f in $deployFiles; do
    echo -e "\n\n\nbegin execute: "$f"\n"
    cd "$(dirname "$f")"
    bash deploy.sh
done
