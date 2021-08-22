#!/bin/bash
# shellcheck disable=SC2086

COMMAND_FILEPATH="$1"
trigger=$(basename "$COMMAND_FILEPATH" .json)
BASHBOT_CONFIG_FILEPATH=${BASHBOT_CONFIG_FILEPATH:-config.json}
tmpconfig=${tmpconfig:-tmp-config.json}
if ! [[ -f "${COMMAND_FILEPATH}" ]]; then
  echo "Missing expected file: ${COMMAND_FILEPATH}"
  exit 1
fi

# echo "COMMAND_FILEPATH:        $COMMAND_FILEPATH"
# echo "BASHBOT_CONFIG_FILEPATH: $BASHBOT_CONFIG_FILEPATH"
# echo "tmpconfig:               $tmpconfig"

# validate the json file

if [[ -n "$(jq '.tools[] | select(.trigger=="'${trigger}'")' ${BASHBOT_CONFIG_FILEPATH})" ]]; then
  echo "Trigger already exists: ${trigger}"
  exit 0
fi
add="$(jq -c '.' ${COMMAND_FILEPATH})"
jq '.tools += ['"${add}"']' ${BASHBOT_CONFIG_FILEPATH} > ${tmpconfig}
mv ${tmpconfig} ${BASHBOT_CONFIG_FILEPATH}

echo "$trigger added to $BASHBOT_CONFIG_FILEPATH"
