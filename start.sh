#!/bin/bash

if [[ -f bashbot.go ]] && [[ -f .env ]]; then
  source .env
  mkdir -p vendor \
    && cd scripts \
    && ./get-config.sh \
    && ./get-vendor-dependencies.sh ../config.json ../vendor \
    && cd ..
  go run bashbot.go
else
  echo "Must run start from project root."
  exit 1
fi