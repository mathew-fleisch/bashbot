FROM golang:1.16 as builder

WORKDIR /bashbot
COPY . .
RUN make go-setup && make go-build

FROM alpine:latest
LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
ENV BASHBOT_CONFIG_FILEPATH=/bashbot/config.json
ENV BASHBOT_ENV_VARS_FILEPATH ""
ENV SLACK_TOKEN ""
ENV LOG_LEVEL "info"
ENV LOG_FORMAT "text"

RUN apk add --update bash curl git make jq \
    && rm /bin/sh && ln -s /bin/bash /bin/sh

WORKDIR /bashbot
COPY --from=builder /bashbot/bin/bashbot-* /usr/local/bin/bashbot
COPY . .
RUN chmod +x /usr/local/bin/bashbot
CMD [ "/bashbot/entrypoint.sh" ]
