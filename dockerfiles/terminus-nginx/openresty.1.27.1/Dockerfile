# syntax=docker/dockerfile:latest

FROM openresty/openresty:1.27.1.2-bookworm-fat

ARG TARGETARCH
ARG ARCH=${TARGETARCH/arm64/aarch64}
ARG ARCH=${ARCH/amd64/x86_64}

LABEL maintainer="majun@terminus.io" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.vendor="Terminus" \
    org.opencontainers.image.authors="hustcer" \
    org.opencontainers.image.title="OpenResty for Terminus" \
    org.opencontainers.image.description="OpenResty for Terminus"

# Use mirrors to speed up installation
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install apt-transport-https ca-certificates locales -y --no-install-recommends --no-install-suggests \
    # Change locale & timezone
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "zh_CN.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen \
    && rm /etc/localtime \
    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'deb https://mirrors.aliyun.com/debian/ bookworm main non-free non-free-firmware contrib' > /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian-security/ bookworm-security main' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian/ bookworm-updates main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.aliyun.com/debian/ bookworm-backports main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian/ bookworm main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian-security/ bookworm-security main' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian/ bookworm-updates main non-free non-free-firmware contrib' >> /etc/apt/sources.list \
    && echo 'deb-src https://mirrors.aliyun.com/debian/ bookworm-backports main non-free non-free-firmware contrib' >> /etc/apt/sources.list

ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8

RUN apt-get install -y --no-install-recommends --no-install-suggests \
       dnsutils tcpdump lsof net-tools telnet nmap ncat \
       vim bat ripgrep \
       wget aria2 unzip \
    && ln -s /usr/bin/batcat /usr/local/bin/bat \
    # Install fd
    && cd /lib \
    && curl -s https://api.github.com/repos/sharkdp/fd/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep ${ARCH}-unknown-linux-musl | aria2c -i - \
    && mkdir fd-latest && tar xvf fd-*.tar.gz --directory=fd-latest \
    && cp -aR fd-latest/**/fd* /usr/local/bin/ \
    && rm -rf fd* \
    && apt-get clean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

# Run the following command to start OpenResty manually
# CMD [ "/usr/local/openresty/bin/openresty", "-g", "daemon off;" ]
