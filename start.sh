#!/bin/bash

source ~/.bashrc
if [[ ! -z $AWS_PUBLIC_SETUP_URL ]]; then
  curl $AWS_PUBLIC_SETUP_URL | bash
fi
if [[ -f bashbot.go ]]; then
  ./install-dependencies.sh
  go run bashbot.go
else
  echo "Must run start from project root."
  exit 1
fi