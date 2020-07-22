#!/bin/bash
STARTING_DIRECTORY=$(pwd)
EXPECTED_COMMANDS="jq aws git"
for expect in $EXPECTED_COMMANDS; do
  if [[ -z $(which "$expect") ]]; then
    echo "Please install $expect to continue"
    exit 1
  fi
done
unset expect
EXPECTED_PRIMARY_VARIABLES="AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY CONFIG_BUCKET"
for expect in $EXPECTED_PRIMARY_VARIABLES; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    exit 1
  fi
done
unset expect
AWS_ACCESS_KEY_ID=$(echo "$AWS_ACCESS_KEY_ID" | tr -d '\n' | tr -d '\r')
AWS_SECRET_ACCESS_KEY=$(echo "$AWS_SECRET_ACCESS_KEY" | tr -d '\n' | tr -d '\r')
CONFIG_BUCKET=$(echo "$CONFIG_BUCKET" | tr -d '\n' | tr -d '\r')
echo "AWS User:"
aws sts get-caller-identity
echo "Get config file: $CONFIG_BUCKET"
aws s3 cp "$CONFIG_BUCKET" .env
echo "Load environment variables..."
# echo "   bucket: '$CONFIG_BUCKET'"
# echo "   awskey: '$AWS_ACCESS_KEY_ID'"
# echo "awssecret: '$AWS_SECRET_ACCESS_KEY'"
if ! [ -f ".env" ]; then
  echo "Missing .env file"
  exit 1
fi
source .env

EXPECTED_VARIABLES="GIT_TOKEN CONFIG_REPO CONFIG_PATH SLACK_TOKEN WELCOME_CHANNEL"
for expect in $EXPECTED_VARIABLES; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    exit 1
  fi
done
unset expect
mkdir -p tmp
cd tmp
CURRENT_PATH=$(pwd)
git clone $CONFIG_REPO
if ! [ -f "$CURRENT_PATH/$CONFIG_PATH/admin.json" ]; then
  echo "Missing admin.json"
  exit 1
fi
if ! [ -f "$CURRENT_PATH/$CONFIG_PATH/config.json" ]; then
  echo "Missing config.json"
  exit 1
fi
if ! [ -f "$CURRENT_PATH/$CONFIG_PATH/messages.json" ]; then
  echo "Missing messages.json"
  exit 1
fi
echo "Copying config from repo: $CONFIG_PATH"
cd $STARTING_DIRECTORY
cp "$CURRENT_PATH/$CONFIG_PATH/admin.json" admin.json
cp "$CURRENT_PATH/$CONFIG_PATH/config.json" config.json
cp "$CURRENT_PATH/$CONFIG_PATH/messages.json" messages.json
rm -rf tmp

mkdir -p vendor
cd vendor

if [[ ! -z $NPM_TOKEN ]]; then
  echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > ~/.npmrc
  echo 'registry=https://registry.npmjs.org/' >> ~/.npmrc
fi

echo "Cloning dependencies"
DEPENDENCIES=$(cat ../config.json | jq -r '.dependencies[] | " \(.install)\(.source)"' | tr '\n' ';')
echo "Install Dependencies:"
echo "$(echo $DEPENDENCIES | tr ';' '\n' | sed -e 's/^\s//g')"
echo "$(eval $(echo $DEPENDENCIES))"

echo "Installing dependencies"
SETUP=$(cat ../config.json | jq -r '.dependencies[].setup' | tr '\n' ';' | sed -e 's/\;/\; /g')
echo "Setting up dependencies:"
echo "$(echo $SETUP | tr ';' '\n' | sed -e 's/^\s//g')"
echo "$(eval $(echo $SETUP))"
