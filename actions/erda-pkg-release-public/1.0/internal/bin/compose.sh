#!/bin/bash

set -o errexit -o nounset -o pipefail
cd "$(dirname "${BASH_SOURCE[0]}")"


# help
if [ "$#" -lt 1 ]; then
    echo "usage:"
    echo "  base compose.sh <dice-version>"
    exit 1
fi

# 校验版本是否存在
version="$1"
if [ ! -d "$PWD/$version" ]; then
    echo "version not exist since 3.8. version: $version"
    exit 1
fi

deploy_yaml="dice.yaml.work"
fdp_deploy_yaml="fdp.yaml.work"

## 生成 dice operator 专有定义部分
cat << END > ${deploy_yaml}
apiVersion: dice.terminus.io/v1beta1
kind: Dice
metadata:
  name: erda
  namespace: {{ .Values.namespace | default "erda-system" }}
spec:
  customDomain:
    openapi: openapi.{{ .Values.domain }}
    ui: {{ .Values.domain }}, *.{{ .Values.domain }}
  cookieDomain: {{ .Values.domain }}
  platformDomain: {{ .Values.domain }}
  diceCluster: {{ .Values.clusterName }}
  size: 'test'
END

## 组合各个 dice.yml 进 dice operator spec 中
echo "  uc:" >> ${deploy_yaml}
cat "$version"/releases/uc/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
echo >> ${deploy_yaml}

#############################################
#############################################
## 注意：新增组件，删除组件需要兼容老版本
## compose 内容写在下面；参照 fdp， pandora

# fdp 单独输出到 fdp.yaml, 用于单独部署 fdp
if [ -f "$version"/releases/fdp/dice.yml ]; then
  echo "  fdp:" >> ${fdp_deploy_yaml}
  cat "$version"/releases/fdp/dice.yml | sed 's/^/    /g' >> ${fdp_deploy_yaml}
  echo >> ${fdp_deploy_yaml}
fi

if [ -f "$version"/releases/elf/dice.yml ]; then
  echo "  elf:" >> ${deploy_yaml}
  cat "$version"/releases/elf/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

# 3.9 delete pandora compoent
if [ -f "$version"/releases/pandora/dice.yml ]; then
  echo "  pandora:" >> ${deploy_yaml}
  cat "$version"/releases/pandora/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

# 3.10
if [ -f "$version"/releases/addon-platform/dice.yml ]; then
  echo "  addonPlatform:" >> ${deploy_yaml}
  cat "$version"/releases/addon-platform/dice.yml | sed 's/^/    /g'  >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

# 3.13
if [ -f "$version"/releases/mesh-controller/dice.yml ]; then
  echo "  meshController:" >> ${deploy_yaml}
  cat "$version"/releases/mesh-controller/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi
if [ -f "$version"/releases/spot-dashboard/dice.yml ]; then
  echo "  spotDashboard:" >> ${deploy_yaml}
  cat "$version"/releases/spot-dashboard/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

## 3.15
if [ -f "$version"/releases/gittar/dice.yml ]; then
  echo "  gittar:" >> ${deploy_yaml}
  cat "$version"/releases/gittar/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

## 3.20
if [ -f "$version"/releases/pmp-backend/dice.yml ]; then
  echo "  pmp:" >> ${deploy_yaml}
  cat "$version"/releases/pmp-backend/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

## 4.0
if [ -f "$version"/releases/hepa/dice.yml ] ; then
  echo "  hepa:" >> ${deploy_yaml}
  cat "$version"/releases/hepa/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
  echo >> ${deploy_yaml}
fi

## erda v1.1.0-rc, remove initJob
if [ -f "$version"/releases/init/dice.yml ]; then
    echo "  initJobs:" >> ${deploy_yaml}
    cat "$version"/releases/init/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

echo "  dice:" >> ${deploy_yaml}
if [ -f "$version"/releases/dice/dice.yml ]; then
    cat "$version"/releases/dice/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
else
    cat "$version"/releases/erda/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

echo "  diceUI:" >> ${deploy_yaml}
if [ -f "$version"/releases/dice-ui/dice.yml ]; then
    cat "$version"/releases/dice-ui/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
else
    cat "$version"/releases/erda-ui/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

echo "  spotAnalyzer:" >> ${deploy_yaml}
if [ -f "$version"/releases/spot-analyzer/dice.yml ]; then
    cat "$version"/releases/spot-analyzer/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
else
    cat "$version"/releases/erda-analyzer/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

echo "  spotTelegraf:" >> ${deploy_yaml}
if [ -f "$version"/releases/spot-telegraf/dice.yml ]; then
    cat "$version"/releases/spot-telegraf/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
else
    cat "$version"/releases/telegraf/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

echo "  spotFilebeat:" >> ${deploy_yaml}
if [ -f "$version"/releases/spot-filebeat/dice.yml ]; then
    cat "$version"/releases/spot-filebeat/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
else
    cat "$version"/releases/beats/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

if [ -f "$version"/releases/tmc/dice.yml ]; then
    echo "  tmc:" >> ${deploy_yaml}
    cat "$version"/releases/tmc/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi
if [ -f "$version"/releases/erda-tmc/dice.yml ]; then
    echo "  tmc:" >> ${deploy_yaml}
    cat "$version"/releases/erda-tmc/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

echo "  spotMonitor:" >> ${deploy_yaml}
if [ -f "$version"/releases/spot-monitor/dice.yml ]; then
    cat "$version"/releases/spot-monitor/dice.yml | sed 's/^/    /g' >> ${deploy_yaml}
    echo >> ${deploy_yaml}
fi

## 新增组件，删除组件需要兼容老版本 end......
#############################################
#############################################


sed -i -e 's/<%/{{/g' ${deploy_yaml}
sed -i -e 's/%>/}}/g' ${deploy_yaml}

sed -i -e 's/<%/{{/g' ${fdp_deploy_yaml}
sed -i -e 's/%>/}}/g' ${fdp_deploy_yaml}


cat ${deploy_yaml} | sed '/  ^$/d' > erda.yaml
rm -f ${deploy_yaml}

cat ${fdp_deploy_yaml} | sed '/  ^$/d' > fdp.yaml
rm -f ${fdp_deploy_yaml}

if [ -f dice.yaml.work-e ]; then
  rm dice.yaml.work-e
fi
