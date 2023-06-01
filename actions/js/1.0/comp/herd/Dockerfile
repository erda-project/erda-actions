FROM {{CENTRAL_REGISTRY}}/erda-x/herd:1.1.20

ARG DICE_VERSION

COPY {{DESTDIR}} .

RUN bootjs=$(node -p "require('./package.json').scripts.start" | \
    sed -n -e 's/^.*herd //p') && \
    bootjs=${bootjs:-'Pampasfile-default.js'} && echo ${bootjs} && \
    npm i @terminus/spot-agent@~${DICE_VERSION} -g && \
    npm link @terminus/spot-agent && \
    spot install -r herd -s ${bootjs} || exit -1;

CMD ["npm", "run", "start"]
