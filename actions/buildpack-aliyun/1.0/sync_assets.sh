#!/usr/bin/env bash

set -exo pipefail

# mkdir for assets
assetsDir=/assets/assets
mkdir -p ${assetsDir}

# sync java-agent
javaAgentDir=${assetsDir}/java-agent
mkdir -p ${javaAgentDir}

function download() {
    first_version=$1
    second_version=$2
    dir=$3
    release_version=${first_version}.${second_version}
    mkdir -p ${dir}/java-agent/${release_version}
    outputFile=${dir}/java-agent/${release_version}/spot-agent.tar.gz
    curl -o ${outputFile} \
    https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/spot/java-agent/action/release/${release_version}/spot-agent.tar.gz
    # check .tar.gz
    printf ${outputFile}
    tar -tzf ${outputFile} >/dev/null
}

## release 3.x
release_3_first_version=3
for i in {15..21}; do
    download ${release_3_first_version} ${i} ${assetsDir}
done

## release 4.x
release_4_first_version=4
for i in {0..0}; do
    download ${release_4_first_version} ${i} ${assetsDir}
done

## release erda 1.x
release_1_first_version=1
for i in {1..2}; do
    download ${release_1_first_version} ${i} ${assetsDir}
done
