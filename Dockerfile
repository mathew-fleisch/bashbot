FROM golang:1.19 as builder

# yq required to parse version from helm chart and inject into build
RUN wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -q -O /usr/bin/yq \
    && chmod +x /usr/bin/yq

WORKDIR /bashbot
COPY . .
RUN make

FROM alpine:latest
LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
ENV BASHBOT_CONFIG_FILEPATH=/bashbot/config.yaml
ENV BASHBOT_ENV_VARS_FILEPATH ""
ENV SLACK_BOT_TOKEN ""
ENV SLACK_APP_TOKEN ""
ENV LOG_LEVEL "info"
ENV LOG_FORMAT "text"
ARG NRUSER=bb

RUN apk add --update --no-cache bash curl git make jq yq \
    && rm -rf /var/cache/apk/* \
    && rm /bin/sh && ln -s /bin/bash /bin/sh \
    && if [[ "${NRUSER}" != "root" ]]; then \
        addgroup -S ${NRUSER} && \
        adduser -D -S ${NRUSER} -G ${NRUSER}; \
    fi

WORKDIR /bashbot
COPY --from=builder --chown=${NRUSER}:${NRUSER} /bashbot/bin/bashbot-* /usr/local/bin/bashbot
COPY . .
RUN chmod +x /usr/local/bin/bashbot \
    && mkdir -p /usr/asdf \
    && chown -R ${NRUSER}:${NRUSER} /usr/asdf \
    && chown -R ${NRUSER}:${NRUSER} /usr/local/bin \
    && chown -R ${NRUSER}:${NRUSER} /bashbot
USER ${NRUSER}
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "which", "bashbot" ]
CMD [ "/bashbot/entrypoint.sh" ]
