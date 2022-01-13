#!/bin/bash

set -o errexit

ts=$(date +%s%3N)
token=${ACTION_ACCESS_TOKEN}
sec=${ACTION_SECRET}
msg=${ACTION_MSG}

if [ -z "${token}" ]; then
  echo "ERROR: access_token is not set"
  exit 1
fi

if [ -z "${sec}" ]; then
  echo "ERROR: secret is not set"
  exit 1
fi

if [ -z "${msg}" ]; then
  echo "ERROR: msg is not set"
  exit 1
fi

sig=$(printf "$ts\n$sec" | openssl dgst -sha256 -hmac "$sec" -binary | base64)

curl -XPOST "https://oapi.dingtalk.com/robot/send?access_token=$token&timestamp=$ts&sign=$sig" -d $msg
