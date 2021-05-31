#!/bin/bash
if [[ -f bashbot.go ]]; then
  # If an aws user+bucket are saved as environment variables, pull the .env file from bucket
  if [[ -n "$AWS_ACCESS_KEY_ID" ]] && [[ -n "$AWS_SECRET_ACCESS_KEY" ]] && [[ -n "$S3_CONFIG_BUCKET" ]]; then
    echo "Attempting to pull .env file from s3 bucket: $(echo $S3_CONFIG_BUCKET | sed -e 's/\//\/ /g')"
    aws sts get-caller-identity
    aws s3 cp ${S3_CONFIG_BUCKET} .env
  fi
  # The .env file doesn't have to come from aws s3 bucket. Verify it exists.
  if [[ -f .env ]]; then
    source .env
    mkdir -p vendor
    pushd scripts
    ./get-config.sh
    ./get-vendor-dependencies.sh ../config.json ../vendor
    popd
    go run bashbot.go
  else
    echo "Must include a .env file at project root. See bashbot read-me for more details."
    exit 1
  fi
else
  echo "Must run start from project root."
  exit 1
fi