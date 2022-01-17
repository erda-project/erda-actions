#!/bin/sh

set -o errexit
set -x

#ts=$(date +%s%3N)
ts=$(date +%s000)
token=${ACTION_ACCESS_TOKEN}
sec=${ACTION_SECRET}
msg_ctx=${ACTION_MSG_CTX}
msg_file=${ACTION_MSG_FILE}

if [ -z "${token}" ]; then
  echo "ERROR: access_token is not set"
  exit 1
fi

if [ -z "${sec}" ]; then
  echo "ERROR: secret is not set"
  exit 1
fi

if [ -z "${msg_ctx}" ]; then
  echo "ERROR: msg_ctx is not set"
  exit 1
fi

if [ -z "${msg_file}" ]; then
  echo "ERROR: msg_file is not set"
  exit 1
fi

jsonnet ${msg_file} --tla-code ctx="${msg_ctx}" -o /tmp/msg.json

echo "INFO: ts: ${ts}"

sig=$(printf "$ts\n$sec" | openssl dgst -sha256 -hmac "$sec" -binary | base64)

curl -XPOST "https://oapi.dingtalk.com/robot/send?access_token=$token&timestamp=$ts&sign=$sig" -H 'Content-Type: application/json' -d@/tmp/msg.json
