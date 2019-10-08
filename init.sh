#!/bin/bash

source .env
cd ./vendor

DEPENDENCIES=$(cat ../config.json | jq -r '.dependencies[] | " \(.install)\(.source)"' | sed -e 's/\$GIT_TOKEN/'"$GIT_TOKEN"'/g' | tr '\n' ';')
echo "Install Dependencies:"
echo "$(echo $DEPENDENCIES | tr ';' '\n' | sed -e 's/^\s//g')"
echo "$(eval $(echo $DEPENDENCIES))"

SETUP=$(cat ../config.json | jq -r '.dependencies[].setup' | tr '\n' ';' | sed -e 's/\;/\; /g')
echo "Setting up dependencies:"
echo "$(echo $SETUP | tr ';' '\n' | sed -e 's/^\s//g')"
echo "$(eval $(echo $SETUP))"