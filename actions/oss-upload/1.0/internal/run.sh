#!/bin/sh
set -e

echo 'oss upload'

cat << EOF > ~/.ossutilconfig
[Credentials]
language=EN
accessKeyID=$ACTION_ACCESS_KEY_ID
accessKeySecret=$ACTION_ACCESS_KEY_SECRET
endpoint=$ACTION_ENDPOINT
stsToken=$ACTION_STS_TOKEN
EOF

ossutil cp "$ACTION_LOCAL_PATH" "oss://$ACTION_BUCKET/$ACTION_OSS_PATH" "--meta=$ACTION_META" -r -f