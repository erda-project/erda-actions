FROM buildpack-deps

ENV NVM_DIR=$HOME
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash \
      && . $HOME/.nvm/nvm.sh \
      && nvm install v12.22.5 \
      && nvm install v14.19.0 \
      && nvm install v16.20.2 \
      && nvm install v18.17.1 \
      && nvm install v20.15.0 \
      && nvm alias default v14.19.0 \
      && nvm use default

COPY actions/js-pack/1.0/index.js /opt/action/run

RUN chmod 777 /opt/action/run
