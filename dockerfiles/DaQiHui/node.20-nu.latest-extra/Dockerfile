# syntax=docker/dockerfile:latest

# Nu: 0.112.2, Node: 20.18.3, npm: 10.8.2, pnpm: 9.15.9, Git: 2.47.3, nvm: 0.40.3
# Bash: 5.2.37, Wget: 1.25.0, aria2c: 1.37.0, curl: 8.14.1, rg: 14.1.1, sd: 1.0.0
# fd: 10.4.2, just: 1.47.1, ossutil: v1.7.19, ast-grep/sg: 0.42.0
# /etc/os-release: Debian GNU/Linux 13.4 (trixie)
# Kernel: Linux ea3c3761f149 6.14.0-1017-azure #17~24.04.1-Ubuntu SMP Mon Dec  1 20:10:50 UTC 2025 x86_64 GNU/Linux
# Build cmd: docker build . -t hustcer/debian-nu:0.112-node.20.18-extra --no-cache
# docker buildx build --platform linux/amd64,linux/arm64 . -t hustcer/debian-nu:0.112-node.20.18-extra --no-cache
# Other tags: terminus/debian-nu:0.112-n20.18-extra,terminus/debian-nu:0.112-extra
FROM hustcer/debian-nu:0.112-node.20.18.lts

ARG TARGETARCH
ARG ARCH=${TARGETARCH/arm64/aarch64}
ARG ARCH=${ARCH/amd64/x86_64}

LABEL org.opencontainers.image.authors="hustcer <hustcer@outlook.com>"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Use mirrors to speed up installation
RUN apt-get update \
    && apt-get upgrade -y \
    # Install sg
    && cd /lib || exit 1 \
    && curl -s https://api.github.com/repos/ast-grep/ast-grep/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep app-${ARCH}-unknown-linux-gnu | aria2c -i - \
    && mkdir sg-latest && unzip app-*.zip -d sg-latest \
    && mv sg-latest/* /usr/local/bin/ \
    # Install fd
    && curl -s https://api.github.com/repos/sharkdp/fd/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep ${ARCH}-unknown-linux-musl | aria2c -i - \
    && mkdir fd-latest && tar xvf fd-*.tar.gz --directory=fd-latest \
    && cp -aR fd-latest/**/fd* /usr/local/bin/ \
    # Install just
    && curl -s https://api.github.com/repos/casey/just/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep ${ARCH}-unknown-linux-musl | aria2c -i - \
    && mkdir just-latest && tar xvf just-*.tar.gz --directory=just-latest \
    && cp -aR just-latest/just* /usr/local/bin/ \
    # Install Ossutils
    && curl https://gosspublic.alicdn.com/ossutil/install.sh | bash \
    && rm -rf /var/lib/apt/lists/* && apt-get autoremove -y \
    && rm -rf just* fd* sg*

CMD [ "nu" ]
