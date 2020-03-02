FROM golang:1.12.6-stretch

LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN mkdir /bashbot

RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt upgrade -y
RUN DEBIAN_FRONTEND=noninteractive apt install -y zip wget iputils-ping curl jq build-essential libssl-dev ssh python python-pip python3 python3-pip openssl file libgcrypt-dev git redis-server sudo build-essential libssl-dev awscli sqlite3 vim

# Get newest version of awscli
RUN pip3 install awscli --upgrade --user

# Install kubectl
RUN echo "Install kubectl"
RUN KUBECTL_VERSION=$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt) && echo "Install kubectl version ${KUBECTL_VERSION}" && curl -sL https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && chmod +x /usr/local/bin/kubectl
# RUN kubectl version

# Install yq
RUN echo "Install yq"
RUN curl -sL https://github.com/mikefarah/yq/releases/download/3.1.0/yq_linux_amd64 -o /usr/local/bin/yq && chmod +x /usr/local/bin/yq
RUN yq --version

# Install kops
RUN KOPS_VERSION=$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4) && echo "Install KOPS: $KOPS_VERSION" && curl -sLO https://github.com/kubernetes/kops/releases/download/$KOPS_VERSION/kops-linux-amd64 && chmod +x kops-linux-amd64 && mv kops-linux-amd64 /usr/local/bin/kops
RUN kops version

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 10.16.3
RUN curl -s https://raw.githubusercontent.com/creationix/nvm/v0.30.1/install.sh | bash \
    && source $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

WORKDIR /bashbot
COPY . .
RUN mkdir -p vendor

RUN cat /bashbot/.env >> ~/.bashrc
RUN source ~/.bashrc
RUN go install -v ./...
RUN go get github.com/nlopes/slack@master

CMD ["/bin/bash", "start.sh"]
