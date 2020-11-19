#!/bin/bash
slack_users_list() {
  IFS='' read -r -d '' usage <<"EOF"
Usage:

slack_users_list slack_api_base slack_token endpoint verbosity output

Note: verbosity is an int bool (0 or 1), output => [json|raw|ids]
EOF
  [ -z "$1" ] && echo "must define slack_api_base url" && echo "$usage" && exit 1;
  SLACK_API_BASE="$1"
  [ -z "$2" ] && echo "must define slack_token" && echo "$usage" && exit 1;
  SLACK_TOKEN="$2"
  [ -z "$3" ] && echo "must define endpoint" && echo "$usage" && exit 1;
  ENDPOINT="$3"
  [ -z "$4" ] && echo "must define verbosity state" && echo "$usage" && exit 1;
  VERBOSE="$4"
  [ -z "$5" ] && echo "must define output display type" && echo "$usage" && exit 1;
  OUTPUT="$5"

  [ $VERBOSE == 1 ] && echo "Get list of users..."
  RESPONSE=$(curl -s "$SLACK_API_BASE/$ENDPOINT?token=$SLACK_TOKEN")
  case "$OUTPUT" in
  "json" )
    echo "$RESPONSE" | jq '.'
    ;;
  "idtags" )
    echo "$RESPONSE" | jq -r '.members[] | select(.deleted == false) | .id' | sort | sed -e 's/\(.*\)/<@\1>/g'
    ;;
  "ids" )
    echo "$RESPONSE" | jq -r '.members[] | select(.deleted == false) | .id' | sort
    ;;
  "raw" )
    echo "$RESPONSE"
    ;;
  *)
    echo "$RESPONSE"
    ;;
esac
}
slack_users_info() {
  IFS='' read -r -d '' usage <<"EOF"
Usage:

slack_users_info slack_api_base slack_token endpoint user_id verbosity output

Note: verbosity is an int bool (0 or 1), output => [json|raw|id]
EOF
  [ -z "$1" ] && echo "must define slack_api_base url" && echo "$usage" && exit 1;
  SLACK_API_BASE="$1"
  [ -z "$2" ] && echo "must define slack_token" && echo "$usage" && exit 1;
  SLACK_TOKEN="$2"
  [ -z "$3" ] && echo "must define endpoint" && echo "$usage" && exit 1;
  ENDPOINT="$3"
  [ -z "$4" ]  && echo "must define user-id" && echo "$usage" && exit 1;
  USER_ID="$4"
  [ -z "$5" ] && echo "must define verbosity state" && echo "$usage" && exit 1;
  VERBOSE="$5"
  [ -z "$6" ] && echo "must define output display type" && echo "$usage" && exit 1;
  OUTPUT="$6"

  [ $VERBOSE == 1 ] && echo "Get user information..."
  RESPONSE=$(curl -s "$SLACK_API_BASE/$ENDPOINT?token=$SLACK_TOKEN&user=$USER_ID")
  case "$OUTPUT" in
  "id" )
    echo "$RESPONSE" | jq -r '.user.id'
    ;;
  "profile" )
    echo "$RESPONSE" | jq -r '.user.profile'
    ;;
  "json" )
    echo "$RESPONSE" | jq '.'
    ;;
  "raw" )
    echo "$RESPONSE"
    ;;
  *)
    echo "$RESPONSE"
    ;;
esac
}