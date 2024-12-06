#FROM {{CENTRAL_REGISTRY}}/erda-x/herd:1.1.20
FROM docker.io/lirbrary/node:{{NODE_VERSION}}
ARG DICE_VERSION

COPY {{DESTDIR}} .

RUN npm -g install herd

RUN bootjs=$(node -p "require('./package.json').scripts.start" | \
    sed -n -e 's/^.*herd //p') && \
    bootjs=${bootjs:-'Pampasfile-default.js'} && echo ${bootjs} && \
    npm i @terminus/spot-agent@4.1.3 && \
    npm link @terminus/spot-agent && \
    spot install -r herd -s ${bootjs} || exit -1;

CMD ["npm", "run", "start"]
