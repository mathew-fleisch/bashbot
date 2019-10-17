#!/bin/bash

source .env
if [[ ! -z $NPM_TOKEN ]]; then
  echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > ~/.npmrc
  echo 'registry=https://registry.npmjs.org/' >> ~/.npmrc
fi
mkdir -p vendor
cd ./vendor

DEPENDENCIES=$(cat ../config.json | jq -r '.dependencies[] | " \(.install)\(.source)"' | sed -e 's/\$GITHUB_TOKEN/'"$GITHUB_TOKEN"'/g' | tr '\n' ';')
echo "Install Dependencies:"
echo "$(echo $DEPENDENCIES | tr ';' '\n' | sed -e 's/^\s//g')"
echo "$(eval $(echo $DEPENDENCIES))"

SETUP=$(cat ../config.json | jq -r '.dependencies[].setup' | tr '\n' ';' | sed -e 's/\;/\; /g')
echo "Setting up dependencies:"
echo "$(echo $SETUP | tr ';' '\n' | sed -e 's/^\s//g')"
echo "$(eval $(echo $SETUP))"