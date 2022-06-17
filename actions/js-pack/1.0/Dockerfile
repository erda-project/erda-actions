FROM buildpack-deps

ENV NVM_DIR=$HOME
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash \
      && . $HOME/.nvm/nvm.sh \
      && nvm install v12.22.5 \
      && nvm install v14.19.0 \
      && nvm alias default v14.19.0 \
      && nvm use default

COPY run /opt/action/run

RUN chmod 777 /opt/action/run
