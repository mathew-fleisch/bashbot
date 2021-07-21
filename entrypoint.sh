#!/bin/bash

if ! command -v bashbot > /dev/null; then
  echo "bashbot is not installed. Please install bashbot and try again."
  exit 1
fi

if ! [ -f "$BASHBOT_CONFIG_FILEPATH" ]; then
  echo "bashbot config file not found. Please create one and try again."
  exit 1
fi

if [ -z "$SLACK_TOKEN" ]; then
  echo "SLACK_TOKEN is not set. Please set it and try again."
  exit 1
fi

# If the log-level doesn't exist, set it to 'info'
LOG_LEVEL=${LOG_LEVEL:-info}
# If the log-format doesn't exist, set it to 'text'
LOG_FORMAT=${LOG_FORMAT:-text}

# Run Bashbot binary passing the config file and the Slack token
bashbot \
  --config-file "$BASHBOT_CONFIG_FILEPATH" \
  --slack-token "$SLACK_TOKEN" \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"
