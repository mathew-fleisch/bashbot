#!/bin/bash

IFS='' read -r -d '' banner <<"EOF"
---------------------------------------
 ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|
---------------------------------------
EOF

IFS='' read -r -d '' help <<"EOF"
Usage: ./bashbot.sh [arguments]

    --action                    [build-ecs,build-docker,pull-configs,push-configs]

      build-ecs --------------> Run circleci job to build docker, and push to ecr (using remote configs) 
          --config-bucket       [s3-bucket]
          --circle-token        [circle-token]
          --circle-project      [circle-organization-fork]

      build-docker -----------> Build the docker container using remote configs 
          --config-bucket       [s3-bucket]

      pull-configs -----------> Make backup of existing configs, pull from bucket
          --config-bucket       [s3-bucket]

      push-configs -----------> Make backup of remote configs, push to bucket
          --config-bucket       [s3-bucket]
EOF
help="$banner\n$help"

[[ ! -f "bashbot.go" ]] && echo "Must execute from project root" && exit 1;
CONFIG_FILES=(.env config.json messages.json admin.json)
PROJECT_ROOT=$(pwd)
TMP_DIR=./tmp

while [[ $# -gt 0 ]] && [[ "$1" == "--"* ]]; do
  opt="$1";
  shift;
  case "$opt" in
      "--" ) break 2;;
      "--action" )
         ACTION="$1"; shift;;
      "--config-bucket" )
         REMOTE_CONFIG_BUCKET="$1"; shift;;
      "--circle-token" )
         CIRCLE_TOKEN="$1"; shift;;
      "--circle-project" )
         CIRCLE_PROJECT="$1"; shift;;
      *) echo >&2 "Invalid option: $@"; exit 1;;
  esac
done

pull_s3_config () {
    [ -z "$1" ] && echo "Missing remote config s3 bucket." && exit 1;
    REMOTE_CONFIG_BUCKET="$1"
    NOW=$(date +%s)
    # Backup local file, pull remote file
    for i in "${CONFIG_FILES[@]}"; do
      mv ./$i $TMP_DIR/local-$NOW-$i
      aws s3 cp $REMOTE_CONFIG_BUCKET/$i $i
      [[ ! -f "$i" ]] && echo "$i file did not download. restoring backup..." && cp $TMP_DIR/local-$NOW-$i $i && exit 1;
    done
    cd $TMP_DIR
    tar -zcvf configs-local-$NOW.tar.gz ./local-$NOW*
    rm -rf ./local-$NOW*
    cd $PROJECT_ROOT
}
push_s3_config () {
    [ -z "$1" ] && echo "Missing remote config s3 bucket." && exit 1;
    REMOTE_CONFIG_BUCKET="$1"
    NOW=$(date +%s)
    # Backup remote file, push local file
    for i in "${CONFIG_FILES[@]}"; do
      aws s3 cp $REMOTE_CONFIG_BUCKET/$i $TMP_DIR/remote-$NOW-$i
      aws s3 cp $i $REMOTE_CONFIG_BUCKET/$i
    done
    cd $TMP_DIR
    tar -zcvf configs-remote-$NOW.tar.gz ./remote-$NOW*
    rm -rf ./remote-$NOW*
    cd $PROJECT_ROOT
}

[ -z "$ACTION" ] && echo "$help" && echo "Must chose an action" && exit 1; 
echo "$banner"
echo "Action: $ACTION"

##
## --action build-ecs
##
if [ "$ACTION" == "build-ecs" ]; then
    [ -z "$REMOTE_CONFIG_BUCKET" ] && echo "Missing remote config s3 bucket." && exit 1;
    [ -z "$CIRCLE_TOKEN" ] && echo "Missing circle token." && exit 1;
    [ -z "$CIRCLE_PROJECT" ] && echo "Missing circle project (organization/fork)." && exit 1;
    pull_s3_config $REMOTE_CONFIG_BUCKET
    CIRCLE_URL="https://circleci.com/gh/$CIRCLE_PROJECT"
    BUILD_URL="https://circleci.com/api/v1.1/project/github/$CIRCLE_PROJECT/tree/master?circle-token=$CIRCLE_TOKEN"
    json=$(jq -c -r -n '{"build_parameters":{"CIRCLE_JOB":"ecs_deploy","REMOTE_CONFIG_BUCKET":"'$REMOTE_CONFIG_BUCKET'"}}')
    response=$(curl -s -X POST --data $json --header "Content-Type:application/json" --url "$BUILD_URL")
    if [ -z "$(echo $response | jq -r 'keys' | grep build_num)" ]; then
      echo "$(echo $response | jq -r '.message') - Note: The circle-token is unique to the project."
      echo "Raw Response:"
      echo "$response"
      exit 1
    fi
    echo "$CIRCLE_URL/$(echo $response | jq -r -c '.build_num')"
    exit 0
fi

##
## --action build-docker
##
if [ "$ACTION" == "build-docker" ]; then
    [ -z "$REMOTE_CONFIG_BUCKET" ] && echo "Missing remote config s3 bucket." && exit 1;
    pull_s3_config $REMOTE_CONFIG_BUCKET
    docker build -t bashbot .
    echo "Run: docker run -it bashbot:latest"
    exit 0
fi

##
## --action pull-configs
##
if [ "$ACTION" == "pull-configs" ]; then
    [ -z "$REMOTE_CONFIG_BUCKET" ] && echo "Missing remote config s3 bucket." && exit 1;
    pull_s3_config $REMOTE_CONFIG_BUCKET
    echo "Local configs backed up and remote config files copied locally"
    exit 0  
fi


##
## --action push-configs
##
if [ "$ACTION" == "push-configs" ]; then
    [ -z "$REMOTE_CONFIG_BUCKET" ] && echo "Missing remote config s3 bucket." && exit 1;
    push_s3_config $REMOTE_CONFIG_BUCKET
    echo "Remote config files backed up and local config files pushed"
    exit 0  
fi


echo "action not found..."
exit 1
  