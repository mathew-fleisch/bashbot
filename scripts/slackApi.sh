#!/bin/bash

IFS='' read -r -d '' banner <<"EOF"
----------------------------------------------
 _____ _            _       ___  ______ _____
/  ___| |          | |     / _ \ | ___ \_   _|
\ `--.| | __ _  ___| | __ / /_\ \| |_/ / | |
 `--. \ |/ _` |/ __| |/ / |  _  ||  __/  | |
/\__/ / | (_| | (__|   <  | | | || |    _| |_
\____/|_|\__,_|\___|_|\_\ \_| |_/\_|    \___/

----------------------------------------------

EOF

IFS='' read -r -d '' help <<"EOF"
Usage: ./slackApi.sh  [arguments]
    --slack-base      [slack-api-base-url]
    --slack-token     [slack-api-key]
    --slack-channel   [slack-channel-id]
    --endpoint        [users.list|users.info|chat.postMessage|chat.postEphemeral]
    --get-id-from-tag [slack-id-tag]
    --get-tag-from-id [slack-id]
    --message         [free-form-text]
    --thread-ts       [timestamp]
    --output          [json|raw]

EOF
PROJECT_ROOT=$(pwd)
[ ! -f "$PROJECT_ROOT/slackApiFunctions.sh" ] && echo "Missing slackApiFunctions.sh file" && exit 1
source "$PROJECT_ROOT/slackApiFunctions.sh"

USER_TAG=""
VERBOSE=0
SLACK_API_BASE=https://slack.com/api
OUTPUT=raw

while [[ $# -gt 0 ]] && [[ "$1" == "--"* ]]; do
  opt="$1";
  shift;
  case "$opt" in
      "--" ) break 2;;
      "--slack-base" )
         SLACK_API_BASE="$1"; shift;;
      "--slack-token" )
         SLACK_TOKEN="$1"; shift;;
      "--slack-channel" )
         CHANNEL="$1"; shift;;
      "--endpoint" )
         ENDPOINT="$1"; shift;;
      "--get-id-from-tag" )
         USER_TAG_TO_ID="$1"; shift;;
      "--get-tag-from-id" )
         USER_ID_TO_TAG="$1"; shift;;
      "--output" )
         OUTPUT="$1"; shift;;
      "--user-id" )
         USER_ID="$1"; shift;;
      "--message" )
         MESSAGE="$@"; shift;;
      "--thread-ts" )
         THREAD_TS="$1"; shift;;
      "--verbose" )
         VERBOSE=1
         ;;
      *) echo >&2 "Invalid option: $@"; exit 1;;
  esac
done

if [ ! -z "$USER_TAG_TO_ID" ]; then
  echo "$USER_TAG_TO_ID" | sed -e 's/<@//g' | sed -e 's/>//g'
  exit 0
fi
if [ ! -z "$USER_ID_TO_TAG" ]; then
  echo "<@$USER_ID_TO_TAG>"
  exit 0
fi

[ -z "$SLACK_TOKEN" ] && echo "must define/export slack token" && echo "$banner$help" && exit 1
[ -z "$SLACK_API_BASE" ] && echo "must define slack base url" && echo "$banner$help" && exit 1
[ -z "$ENDPOINT" ] && echo "must define endpoint" && echo "$banner$help" && exit 1

case "$ENDPOINT" in
  "users.list" )
    slack_users_list $SLACK_API_BASE $SLACK_TOKEN $ENDPOINT $VERBOSE $OUTPUT
    ;;
  "users.info" )
    slack_users_info $SLACK_API_BASE $SLACK_TOKEN $ENDPOINT $USER_ID $VERBOSE $OUTPUT
    ;;
  "chat.postMessage" )
    curl -s -X POST -H 'Content-Type: application/json' -H 'Authorization: Bearer '$SLACK_TOKEN -d '{"text": "'"
    $MESSAGE"'", "channel": "'$CHANNEL'", "thread_ts": "'$THREAD_TS'"}' $SLACK_API_BASE/$ENDPOINT
    ;;
  "chat.postEphemeral" )
    curl -s -X POST -H 'Content-Type: application/json' -H 'Authorization: Bearer '$SLACK_TOKEN -d '{"text": "'"$MESSAGE"'", "channel": "'$CHANNEL'", "user": "'$USER_ID'"}' $SLACK_API_BASE/$ENDPOINT
    ;;
  *) echo >&2 "endpoint not found: $@"; exit 1;;
esac