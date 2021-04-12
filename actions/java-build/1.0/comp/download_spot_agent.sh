#!/usr/bin/env bash

set -exo pipefail

versions=(
3.14
3.15
3.16
3.17
)

# mkdir for assets
compDir=/opt/action/comp
mkdir -p ${compDir}

# download spot-agent
spotAgentDir=${compDir}/spot-agent
mkdir -p ${spotAgentDir}
for v in ${versions[@]}; do
    mkdir -p ${spotAgentDir}/${v}
    outputFile=${spotAgentDir}/${v}/spot-agent.tar.gz
    curl -o ${outputFile} \
    https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/spot/java-agent/action/release/${v}/spot-agent.tar.gz
    # check .tar.gz
    tar -xzf ${outputFile} -C /opt/action/comp/
done
