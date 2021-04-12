#!/usr/bin/env bash

set -exo pipefail

fonts=(
simhei.ttf
simsun.ttc
yahei.ttf
)

fontDir=/opt/action/comp/fonts
mkdir -p ${fontDir}

for font in ${fonts[@]}; do
    outputFile=${fontDir}/${font}
    curl -o ${outputFile} \
    http://terminus-dice.oss-cn-hangzhou.aliyuncs.com/devops/fonts/${font}
done