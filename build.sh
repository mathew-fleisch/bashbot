#!/bin/bash

IFS='' read -r -d '' help <<"EOF"
---------------------------------------
 ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|
---------------------------------------
Usage: ./build.sh [arguments]

    --type           [ecs,docker]
      type[ecs]
          --config-bucket       [s3-bucket]
          --circle-token        [circle-token]
          --circle-project      [circle-organization-fork]
      type[docker]
          --config-bucket       [s3-bucket]
EOF

while [[ $# -gt 0 ]] && [[ "$1" == "--"* ]]; do
  opt="$1";
  shift;
  case "$opt" in
      "--" ) break 2;;
      "--type" )
         BUILD_TYPE="$1"; shift;;
      "--config-bucket" )
         REMOTE_CONFIG_BUCKET="$1"; shift;;
      "--circle-token" )
         CIRCLE_TOKEN="$1"; shift;;
      "--circle-project" )
         CIRCLE_PROJECT="$1"; shift;;
      *) echo >&2 "Invalid option: $@"; exit 1;;
  esac
done

get_s3_config () {
    [ -z "$1" ] && echo "Missing remote config s3 bucket." && exit 1;
    REMOTE_CONFIG_BUCKET="$1"
    aws s3 cp $REMOTE_CONFIG_BUCKET/.env .env
    [[ ! -f ".env" ]] && echo ".env file did not download" && exit 1;
    aws s3 cp $REMOTE_CONFIG_BUCKET/config.json config.json
    [[ ! -f "config.json" ]] && echo "config.json file did not download" && exit 1;
    aws s3 cp $REMOTE_CONFIG_BUCKET/messages.json messages.json
    [[ ! -f "messages.json" ]] && echo "messages.json file did not download" && exit 1;
    aws s3 cp $REMOTE_CONFIG_BUCKET/admin.json admin.json
    [[ ! -f "admin.json" ]] && echo "admin.json file did not download" && exit 1;
}

[[ ! -f "bashbot.go" ]] && echo "Must execute from project root" && exit 1;
[ -z "$BUILD_TYPE" ] && echo "$help" && echo "Must chose a build type" && exit 1; 
echo "Building: $BUILD_TYPE"
if [ "$BUILD_TYPE" == "ecs" ]; then
    [ -z "$REMOTE_CONFIG_BUCKET" ] && echo "Missing remote config s3 bucket." && exit 1;
    [ -z "$CIRCLE_TOKEN" ] && echo "Missing circle token." && exit 1;
    [ -z "$CIRCLE_PROJECT" ] && echo "Missing circle project (organization/fork)." && exit 1;
    get_s3_config $REMOTE_CONFIG_BUCKET
    CIRCLE_URL="https://circleci.com/gh/$CIRCLE_PROJECT"
    BUILD_URL="https://circleci.com/api/v1.1/project/github/$CIRCLE_PROJECT/tree/master?circle-token=$CIRCLE_TOKEN"
    json=$(jq -c -r -n '{"build_parameters":{"CIRCLE_JOB":"ecs_deploy","REMOTE_CONFIG_BUCKET":"'$REMOTE_CONFIG_BUCKET'"}}')
    response=$(curl -s -X POST --data $json --header "Content-Type:application/json" --url "$BUILD_URL")
    echo "$CIRCLE_URL/$(echo $response | jq -r -c '.build_num')"
    exit 0
fi


if [ "$BUILD_TYPE" == "docker" ]; then
    [ -z "$REMOTE_CONFIG_BUCKET" ] && echo "Missing remote config s3 bucket." && exit 1;
    get_s3_config $REMOTE_CONFIG_BUCKET
    docker build -t bashbot .
    echo "Run: docker run -it bashbot:latest"
    exit 0
fi
echo "Build type not found..."
exit 1
  