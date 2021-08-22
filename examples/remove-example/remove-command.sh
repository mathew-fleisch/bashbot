#!/bin/bash
# shellcheck disable=SC2086

trigger="$1"
BASHBOT_CONFIG_FILEPATH=${BASHBOT_CONFIG_FILEPATH:-config.json}
tmpconfig=${tmpconfig:-tmp-config.json}


# echo "BASHBOT_CONFIG_FILEPATH: $BASHBOT_CONFIG_FILEPATH"
# echo "tmpconfig:               $tmpconfig"

# validate the json file

if [[ -z "$(jq '.tools[] | select(.trigger=="'${trigger}'")' ${BASHBOT_CONFIG_FILEPATH})" ]]; then
  echo "Trigger not found to remove: ${trigger}"
  exit 0
fi


jq 'del(.tools[] | select(.trigger=="'${trigger}'"))' ${BASHBOT_CONFIG_FILEPATH} > ${tmpconfig}
mv ${tmpconfig} ${BASHBOT_CONFIG_FILEPATH}

echo "$trigger removed from $BASHBOT_CONFIG_FILEPATH"



