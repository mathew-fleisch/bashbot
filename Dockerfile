FROM ubuntu:20.04

LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
WORKDIR /root
RUN apt update \
    && DEBIAN_FRONTEND=noninteractive apt dist-upgrade -y \
    && DEBIAN_FRONTEND=noninteractive apt install -y curl sudo golang \
    && rm -rf /var/lib/apt/lists/* \
    && echo "%sudo ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
RUN curl -s https://s3.amazonaws.com/download.draios.com/dependencies/get-dependency-installer.sh | bash
WORKDIR /root/dependency-installer
COPY dependencies.yaml dependencies.yaml
RUN ./bootstrap-build.sh \
    && ./installer.sh ./dependencies.yaml

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 10.16.3
RUN curl -s https://raw.githubusercontent.com/creationix/nvm/v0.30.1/install.sh | bash \
    && source $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

RUN mkdir -p /bashbot
WORKDIR /bashbot
COPY . .
RUN mkdir -p vendor
RUN go install -v ./...
RUN go get github.com/nlopes/slack@master

CMD ["/bin/bash", "start.sh"]
