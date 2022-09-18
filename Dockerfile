FROM golang:1.19 as builder

WORKDIR /bashbot
COPY . .
RUN make

FROM alpine:latest
LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
ENV BASHBOT_CONFIG_FILEPATH=/bashbot/config.json
ENV BASHBOT_ENV_VARS_FILEPATH ""
ENV SLACK_BOT_TOKEN ""
ENV SLACK_APP_TOKEN ""
ENV LOG_LEVEL "info"
ENV LOG_FORMAT "text"
ARG NRUSER=bb

RUN apk add --update --no-cache bash curl git make jq \
    && rm -rf /var/cache/apk/* \
    && addgroup -S ${NRUSER} \
    && adduser -D -S ${NRUSER} -G ${NRUSER} \
    && rm /bin/sh && ln -s /bin/bash /bin/sh

WORKDIR /bashbot
COPY --from=builder --chown=${NRUSER}:${NRUSER} /bashbot/bin/bashbot-* /usr/local/bin/bashbot
COPY . .
RUN chmod +x /usr/local/bin/bashbot \
    && mkdir -p /usr/asdf \
    && chown -R ${NRUSER}:${NRUSER} /usr/asdf \
    && chown -R ${NRUSER}:${NRUSER} /bashbot
USER ${NRUSER}

CMD [ "/bashbot/entrypoint.sh" ]
