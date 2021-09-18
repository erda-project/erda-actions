# Git: git version 2.30.2
# gcc version 10.2.1 20210110 (Debian 10.2.1-6)
# Node: 14.17.5 LTS, npm: 6.14.14, yarn: 1.22.5, herd: 1.1.11
# /etc/os-release: Debian: Debian GNU/Linux 11 (bullseye)
# Kernel: Linux 08f4f27793e7 5.10.47-linuxkit Sat Jul 3 21:51:47 2021 x86_64
# Build cmd: docker build . -t terminus/debian-node.14-herd.1.1.11
# Other tags: terminus/debian-herd:1.1.11-n14.17,terminus/debian-herd:1.1.11
FROM node:14.17.5-bullseye-slim

LABEL maintainer=hustcer<majun@terminus.io>

# Use mirrors to speed up installation
RUN cp /etc/apt/sources.list /etc/apt/sources.list.bak \
    && apt update \
    && apt upgrade -y \
    && apt install apt-transport-https ca-certificates locales -y --no-install-recommends --no-install-suggests \
    && echo 'deb https://mirrors.aliyun.com/debian/ bullseye main contrib non-free' > /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian/ bullseye-updates main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian/ bullseye-backports main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian-security bullseye-security main contrib non-free' >> /etc/apt/sources.list \
    # Change locale & timezone
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "zh_CN.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen \
    # Fix: copy and paste CJK characters in terminal.
    && echo 'set input-meta on' > /root/.inputrc \
    && echo 'set output-meta on' > /root/.inputrc \
    && echo 'set convert-meta off' > /root/.inputrc \
    && echo 'set enable-meta-key on' > /root/.inputrc \
    && echo "Asia/Shanghai" > /etc/timezone \
    && cp -a /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8
ENV NODEJS_ORG_MIRROR="https://npm.taobao.org/dist"
ENV SASS_BINARY_SITE="https://npm.taobao.org/mirrors/node-sass"

RUN npm config set @terminus:registry https://registry.npm.terminus.io/ \
    && npm config set registry https://registry.npm.taobao.org/ \
    && apt update \
    && apt install -y --no-install-recommends --no-install-suggests \
       # Network related utils
       wget curl telnet tcpdump net-tools dnsutils iputils-ping \
       # Dev related
       git gcc make g++ vim \
       # SEE: https://github.com/BurntSushi/ripgrep
       # SEE: https://github.com/sharkdp/fd#installation
       tree fd-find ripgrep htop bzip2 unzip lsof \
    # Make /bin/sh symlink to bash instead of dash:
    && echo "dash dash/sh boolean false" | debconf-set-selections \
    && DEBIAN_FRONTEND=noninteractive dpkg-reconfigure dash \
    # fnm --- Reserved expansion point:
    # See: https://github.com/Schniz/fnm, Fast and simple Node version manager
    # Make it easier and faster to switch to other node versions in the future
    # Just do: fnm install 16 && fnm use 16, will finish in 10 seconds.
    && curl -fsSL https://fnm.vercel.app/install | bash \
    # Fix: fnm can't infer your shell based on the process tree.
    && sed -i 's/fnm env/fnm env --shell=bash/' /root/.bashrc \
    # Setup root alias
    && echo "# Setup user alias" >> /root/.bashrc \
    && echo "alias ..='cd ..'" >> /root/.bashrc \
    && echo "alias fd='fdfind'" >> /root/.bashrc \
    && echo "alias ls='ls --color'" >> /root/.bashrc \
    && echo "alias ll='ls -l --color'" >> /root/.bashrc \
    # Fix vim encoding for cjk characters
    && echo 'set enc=utf8' >> /etc/vim/vimrc \
    && echo 'set fencs=utf8,gbk,gb2312,gb18030' >> /etc/vim/vimrc \
    # Add --unsafe-perm to fix 'Error: EACCES: permission denied, mkdir...' issue while install node-sass
    && npm i -g --unsafe-perm @terminus/herd@1.1.11 \
    && rm -rf /var/lib/apt/lists/* && apt autoremove -y \
    # Add dice and erda user and create home dir
    && useradd -m -s /bin/bash dice \
    && useradd -m -s /bin/bash erda

CMD [ "herd" ]
