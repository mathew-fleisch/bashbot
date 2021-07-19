FROM mathewfleisch/tools:latest
LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"
ENV AWS_ACCESS_KEY_ID ""
ENV AWS_SECRET_ACCESS_KEY ""
ENV S3_CONFIG_BUCKET ""
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
RUN go install -v ./...
RUN go get github.com/slack-go/slack@master

CMD /bin/sh -c ". ${ASDF_DATA_DIR}/asdf.sh && ./entrypoint.sh"