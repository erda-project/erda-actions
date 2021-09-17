
# Git: git version 2.30.2
# gcc: 10.3.1 20210424
# Node: 14.17.5 LTS, npm: 6.14.14, yarn: 1.22.5, herd: 1.1.11
# /etc/os-release: Alpine Linux v3.14
# Kernel(uname -a): Linux 55f39d8cbcc0 5.10.47-linuxkit Sat Jul 3 21:51:47 2021 x86_64
# Build cmd: docker build . -t terminus/alpine-node.14-herd.1.1.11
# Other tags: terminus/alpine-herd:1.1.11-n14.17,terminus/alpine-herd:1.1.11
# Notice: 由于Alpine使用musl libc而不是glibc，而大多数流行的发行版使用后者，
#         为了避免兼容性问题现阶段不建议在Alpine里面构建或者部署Node应用
FROM node:14.17.5-alpine3.14

LABEL maintainer=hustcer<majun@terminus.io>

ENV LANG=en_US.UTF-8
ENV NODEJS_ORG_MIRROR="https://npm.taobao.org/dist"
ENV SASS_BINARY_SITE="https://npm.taobao.org/mirrors/node-sass"

# Use aliyun repo to speed up installation
# Testing repo: https://mirrors.aliyun.com/alpine/edge/testing/
RUN echo "https://mirrors.aliyun.com/alpine/v3.14/main/" > /etc/apk/repositories \
    && echo "https://mirrors.aliyun.com/alpine/v3.14/community/" >> /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# Use terminus npm registry
RUN npm config set @terminus:registry https://registry.npm.terminus.io/ \
    && npm config set registry https://registry.npm.taobao.org/ \
    && apk add --no-cache busybox-extras bash \
       # Dev related: install libstdc++ libc6-compat, etc. to fix node-canvas related issue
       # build-base will install: binutils, gcc, g++, make, patch, etc.
       build-base libstdc++ libc6-compat cairo-dev jpeg-dev pango-dev giflib-dev \
       vim git git-subtree python2 \
       # Fix: Error loading shared library libuuid.so.1
       libuuid \
       # Fix: Fontconfig error: Cannot load default config file
       fontconfig \
       # Network related utils
       curl tcpdump ca-certificates \
       # SEE: https://github.com/BurntSushi/ripgrep
       # SEE: https://github.com/sharkdp/fd#installation
       tree fd htop ripgrep \
    # Fix: Error loading shared library ld-linux-x86-64.so.2: No such file or directory
    && ln -s /lib/libc.musl-x86_64.so.1 /lib/ld-linux-x86-64.so.2 \
    # Fix: Error loading shared library libresolv.so.2
    && ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2 \
    # Change default shell from ash to bash
    && sed -i -e "s/bin\/ash/bin\/bash/" /etc/passwd \
    # Setup root alias
    && echo "alias vi='vim'" >> /root/.bashrc \
    && echo "alias ..='cd ..'" >> /root/.bashrc \
    && echo "alias ll='ls -l'" >> /root/.bashrc \
    # Add --unsafe-perm to fix 'Error: EACCES: permission denied, mkdir...' issue while install node-sass
    && npm i -g --unsafe-perm @terminus/herd@1.1.11 \
    && rm -rf /var/cache/apk/* \
    && echo "" > /root/.ash_history \
    # Add dice and erda user and create home dir, Don't assign a password
    && adduser -h dice -Ds /bin/bash dice \
    && adduser -h erda -Ds /bin/bash erda

CMD [ "herd" ]
