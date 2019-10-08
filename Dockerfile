FROM golang:1.12.6-stretch

LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN mkdir /bashbot

RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt upgrade -y
RUN DEBIAN_FRONTEND=noninteractive apt install -y zip wget iputils-ping curl jq build-essential libssl-dev ssh python python-pip python3 python3-pip openssl file libgcrypt-dev git redis-server sudo build-essential libssl-dev awscli vim

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 10.16.3
RUN curl https://raw.githubusercontent.com/creationix/nvm/v0.30.1/install.sh | bash \
    && source $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

WORKDIR /bashbot
COPY . .

RUN cat /bashbot/.env >> ~/.bashrc
RUN source ~/.bashrc
RUN eval $(cat /bashbot/.env) && echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > ~/.npmrc
RUN echo 'registry=https://registry.npmjs.org/' >> ~/.npmrc

WORKDIR /bashbot
RUN go install -v ./...


CMD ["/bin/bash", "start.sh"]
