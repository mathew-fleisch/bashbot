FROM alpine:latest
LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
ENV BASHBOT_CONFIG_FILEPATH=/bashbot/config.json
ENV BASHBOT_ENV_VARS_FILEPATH ""
ENV SLACK_TOKEN ""
ENV LOG_LEVEL "info"
ENV LOG_FORMAT "text"
ENV ASDF_DATA_DIR /root/.asdf

RUN apk add --update bash curl git make go jq docker python3 py3-pip openssh vim \
    && rm /bin/sh && ln -s /bin/bash /bin/sh \
    && ln -s /usr/bin/python3 /usr/local/bin/python

# Install asdf dependencies
WORKDIR /root
COPY .tool-versions /root/.tool-versions
COPY pin /root/pin
RUN mkdir -p $ASDF_DATA_DIR \
    && git clone --depth 1 https://github.com/asdf-vm/asdf.git $ASDF_DATA_DIR \
    && . $ASDF_DATA_DIR/asdf.sh \
    && echo -e '\n. $ASDF_DATA_DIR/asdf.sh' >> $HOME/.bashrc \
    && echo -e '\n. $ASDF_DATA_DIR/asdf.sh' >> $HOME/.profile \
    && asdf update \
    && while IFS= read -r line; do asdf plugin add $(echo "$line" | awk '{print $1}'); done < .tool-versions \
    && asdf install

RUN mkdir -p /bashbot
WORKDIR /bashbot
COPY . .
RUN mkdir -p vendor
RUN . ${ASDF_DATA_DIR}/asdf.sh \
    && make build \
    && mv bin/bashbot-* /usr/local/bin/bashbot \
    && chmod +x /usr/local/bin/bashbot

CMD /bin/sh -c ". ${ASDF_DATA_DIR}/asdf.sh && ./entrypoint.sh"
