docker build -t erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/android-gradle-node:v29 --no-cache \
 --build-arg http_proxy=http://30.43.41.107:1087  --build-arg https_proxy=http://30.43.41.107:1087 .