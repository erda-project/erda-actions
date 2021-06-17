
FROM registry.erda.cloud/erda/terminus-nodejs:node-9.11.1-npm-5.8.0
ARG FORCE_DEP="true"
ARG FORCE_UPDATE_SNAPSHOT="false"
ARG PACKAGE_LOCK_DIR
COPY ${PACKAGE_LOCK_DIR} /app
ARG NODE_OPTIONS
ARG DEP_CMD="echo hello && echo I'm linjun && echo I\"\""\""m linjun"
RUN cd /app && eval NODE_OPTIONS=${NODE_OPTIONS} ${DEP_CMD}
COPY /src /app
WORKDIR /app
RUN NODE_OPTIONS=${NODE_OPTIONS} npm run build
