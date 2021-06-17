FROM {{CENTRAL_REGISTRY}}/erda-actions/terminus-herd:1.1.8-node12

ARG DICE_VERSION

RUN echo "Asia/Shanghai" | tee /etc/timezone

COPY {{DESTDIR}} .

# Set special timezone
RUN bootjs=$(node -p "require('./package.json').scripts.start" | \
    sed -n -e 's/^.*herd //p') && \
    bootjs=${bootjs:-'Pampasfile-default.js'} && echo ${bootjs} && \
    npm i @terminus/spot-agent@~${DICE_VERSION} -g && \
    npm link @terminus/spot-agent && \
    spot install -r herd -s ${bootjs} || exit -1;

CMD ["npm", "run", "start"]
