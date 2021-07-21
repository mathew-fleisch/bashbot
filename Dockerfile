FROM mathewfleisch/tools:latest
# See this repo for the parent Dockerfile: https://github.com/mathew-fleisch/tools
LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
ENV BASHBOT_CONFIG_FILEPATH=/bashbot/config.json
ENV SLACK_TOKEN ""
ENV LOG_LEVEL "info"
ENV LOG_FORMAT "text"
ENV ASDF_DATA_DIR /opt/asdf

USER root
# Install asdf dependencies
WORKDIR /root
COPY .tool-versions /root/.tool-versions
COPY pin /root/pin
RUN . ${ASDF_DATA_DIR}/asdf.sh  \
    && asdf update \
    && while IFS= read -r line; do asdf plugin add $(echo "$line" | awk '{print $1}'); done < .tool-versions \
    && asdf install

RUN mkdir -p /bashbot
WORKDIR /bashbot
COPY . .
RUN mkdir -p vendor
RUN . ${ASDF_DATA_DIR}/asdf.sh \
    && make setup \
    && make build \
    && mv bin/bashbot-* /usr/local/bin/bashbot

CMD /bin/sh -c ". ${ASDF_DATA_DIR}/asdf.sh && ./entrypoint.sh"