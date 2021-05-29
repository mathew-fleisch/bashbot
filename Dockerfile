FROM ubuntu:20.04

LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
ENV AWS_ACCESS_KEY_ID ""
ENV AWS_SECRET_ACCESS_KEY ""
ENV S3_CONFIG_BUCKET ""
COPY scripts/.tool-versions /root/.
ENV ASDF_DATA_DIR=/opt/asdf

WORKDIR /root
# Apt dependencies
RUN apt update \
    && DEBIAN_FRONTEND=noninteractive apt dist-upgrade -y \
    && DEBIAN_FRONTEND=noninteractive apt install -y curl wget apt-utils python3 python3-pip make build-essential openssl lsb-release libssl-dev apt-transport-https ca-certificates iputils-ping git vim jq zip sudo golang ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && echo "%sudo ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
# asdf dependencies
RUN mkdir -p $ASDF_DATA_DIR \
    && chmod -R g+w $ASDF_DATA_DIR \
    && git clone https://github.com/asdf-vm/asdf.git ${ASDF_DATA_DIR} --branch v0.8.0 \
    && echo "export ASDF_DATA_DIR=${ASDF_DATA_DIR}" | tee -a /root/.bashrc \
    && echo ". ${ASDF_DATA_DIR}/asdf.sh" | tee -a /root/.bashrc \
    && . ${ASDF_DATA_DIR}/asdf.sh  \
    && asdf plugin add awscli \
    && asdf plugin add golang \
    && asdf plugin add helm \
    && asdf plugin add helmfile \
    && asdf plugin add k9s \
    && asdf plugin add kubectl \
    && asdf plugin add kubectx \
    && asdf plugin add shellcheck \
    && asdf plugin add terraform \
    && asdf plugin add terragrunt \
    && asdf plugin add tflint \
    && asdf plugin add yq \
    && asdf install
RUN mkdir -p /bashbot
WORKDIR /bashbot
COPY . .
RUN mkdir -p vendor
RUN go install -v ./...
RUN go get github.com/slack-go/slack@master

CMD /bin/sh -c ". ${ASDF_DATA_DIR}/asdf.sh && ./entrypoint.sh"