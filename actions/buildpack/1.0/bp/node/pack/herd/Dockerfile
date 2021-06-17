FROM {{BP_DOCKER_BASE_REGISTRY}}/erda-actions/terminus-herd:1.1.9-beta.1

COPY /bp/pack/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /app

ADD /app /app

ARG DICE_VERSION

RUN \
    bootjs=$(node -p "require('./package.json').scripts.start" | \
    sed -n -e 's/^.*herd //p') && \
    bootjs=${bootjs:-'Pampasfile-default.js'} && echo ${bootjs} && \
    npm i @terminus/spot-agent@~${DICE_VERSION} -g && \
    npm link @terminus/spot-agent && \
    spot install -r herd -s ${bootjs} || exit -1;

ENTRYPOINT ["/entrypoint.sh"]
