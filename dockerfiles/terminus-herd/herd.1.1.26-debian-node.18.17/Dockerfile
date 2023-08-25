# Git: git version 2.39.2, VIM - Vi IMproved 9.0
# gcc version 12.2.0, curl: 7.88.1, Wget: 1.21.3, GNU bash: 5.2.15
# Node: 18.17.1 LTS, npm: 9.6.7, yarn: 1.22.19, herd: 1.1.26, pnpm: 8.6.12, fnm: 1.35.1
# /etc/os-release: Debian GNU/Linux 12 (bookworm)
# Kernel: Linux 674a7d667c9d 5.15.49-linuxkit-pr #1 SMP Thu May 25 07:17:40 UTC 2023 x86_64 GNU/Linux
# Build cmd: docker build . -t terminus/debian-node.18-herd.1.1.26
# Other tags: terminus/debian-herd:1.1.26-n18.17,terminus/debian-herd:1.1.26
FROM node:18.17.1-bookworm-slim

LABEL maintainer=hustcer<majun@terminus.io>

# Use mirrors to speed up installation
RUN apt update \
    && apt upgrade -y \
    && apt install apt-transport-https ca-certificates locales -y --no-install-recommends --no-install-suggests \
    && echo 'deb https://mirrors.aliyun.com/debian/ bookworm main non-free non-free-firmware contrib' > /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian-security/ bookworm-security main' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian/ bookworm-updates main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian/ bookworm-backports main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian/ bookworm main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian-security/ bookworm-security main' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian/ bookworm-updates main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian/ bookworm-backports main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
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
ENV NODEJS_ORG_MIRROR="https://npmmirror.com/dist"
ENV SASS_BINARY_SITE="https://npmmirror.com/mirrors/node-sass"

RUN npm config set @terminus:registry https://registry.npm.terminus.io/ \
    && npm config set registry https://registry.npmmirror.com/ \
    && apt update \
    && apt install -y --no-install-recommends --no-install-suggests \
       # Network related utils
       wget curl telnet tcpdump net-tools dnsutils iputils-ping \
       # Dev related
       git gcc make g++ vim \
       # SEE: https://github.com/BurntSushi/ripgrep
       # SEE: https://github.com/sharkdp/fd#installation
       tree fd-find ripgrep htop bzip2 unzip lsof procps \
    # Make /bin/sh symlink to bash instead of dash:
    # REF: https://groups.google.com/g/linux.debian.user/c/2CyfM5UrSOI?pli=1
    && dpkg-divert --remove --no-rename /usr/share/man/man1/sh.1.gz \
    && dpkg-divert --remove --no-rename /bin/sh \
    && ln -sf bash.1.gz /usr/share/man/man1/sh.1.gz \
    && ln -sf bash /bin/sh \
    && dpkg-divert --add --local --no-rename /usr/share/man/man1/sh.1.gz \
    && dpkg-divert --add --local --no-rename /bin/sh \
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
    && npm i -g --unsafe-perm @terminus/herd@1.1.26 \
    && npm i -g pnpm@latest \
    # Fix "Unable to find the global bin directory" error while running `pnpm i -g`
    && pnpm config set global-bin-dir '/usr/local/bin' \
    && rm -rf /var/lib/apt/lists/* && apt autoremove -y \
    # Add erda user and create home dir
    && useradd -m -s /bin/bash erda

CMD [ "herd" ]
