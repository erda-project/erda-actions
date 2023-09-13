FROM registry.erda.cloud/erda-x/herd:1.1.20

COPY /bp/pack/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /app

ADD /app /app

RUN \
    bootjs=$(node -p "require('./package.json').scripts.start" | \
    sed -n -e 's/^.*herd //p') && \
    bootjs=${bootjs:-'Pampasfile-default.js'} && echo ${bootjs} && \
    npm i @terminus/spot-agent@~2.5.0 -g && \
    npm link @terminus/spot-agent && \
    spot install -r herd -s ${bootjs} || exit -1;

ENTRYPOINT ["/entrypoint.sh"]
