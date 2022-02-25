#!/bin/bash

# First check if bashbot is installed and pull that version
if command -v bashbot > /dev/null; then 
    command -v bashbot
    bashbot --version
    exit 0
fi

# Next check if ./bin/bashbot-${os}-${arch} exists and pull that version
arch=amd64
[ "$(uname -m)" == "aarch64" ] && arch=arm64
os=$(uname | tr '[:upper:]' '[:lower:]')
if [ -f "./bin/bashbot-${os}-${arch}" ]; then
    echo "./bin/bashbot-${os}-${arch} --version"
    ./bin/bashbot-${os}-${arch} --version
    exit 0
fi

# Finally, check if bashbot go source exists and pull that version
go_filename=cmd/bashbot/bashbot.go
if [[ -f "./${go_filename}" ]] && [[ -f "./Makefile" ]]; then
    make go-setup
    make go-version
    exit 0
fi

echo "Could not determine the current version of bashbot"
exit 0