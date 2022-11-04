#!/bin/bash
# shellcheck disable=SC1090

echo "
 ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|"

if ! command -v bashbot > /dev/null; then
  echo "bashbot is not installed. Please install bashbot and try again."
  exit 1
fi

# If .env file is present, load it.
if [ -f "$BASHBOT_ENV_VARS_FILEPATH" ]; then
  . "$BASHBOT_ENV_VARS_FILEPATH"
fi

if ! [ -f "$BASHBOT_CONFIG_FILEPATH" ]; then
  echo "bashbot config file not found. Please create one and try again."
  exit 1
fi

if [ -z "$SLACK_BOT_TOKEN" ]; then
  echo "SLACK_BOT_TOKEN is not set. Please set it and try again."
  exit 1
fi

if [ -z "$SLACK_APP_TOKEN" ]; then
  echo "SLACK_APP_TOKEN is not set. Please set it and try again."
  exit 1
fi
mkdir -p vendor

# If the log-level doesn't exist, set it to 'info'
LOG_LEVEL=${LOG_LEVEL:-info}
# If the log-format doesn't exist, set it to 'text'
LOG_FORMAT=${LOG_FORMAT:-text}

# Run install-dependencies path
bashbot install-dependencies \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"

# Run Bashbot binary passing the config file and the Slack token
bashbot run \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"
