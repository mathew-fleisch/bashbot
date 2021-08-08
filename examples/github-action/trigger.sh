#!/bin/bash

github_base="${github_base:-api.github.com}"
expected_variables="BASHBOT_CONFIG_FILEPATH SLACK_TOKEN SLACK_CHANNEL SLACK_USERID REPO_OWNER REPO_NAME GITHUB_TOKEN"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 1
  fi
done
headers="-sH \"Accept: application/vnd.github.everest-preview+json\" -H \"Authorization: token ${GITHUB_TOKEN}\" -X POST"
curl ${headers} \
  --data '{"event_type":"trigger-github-action","client_payload": {"channel":"'${SLACK_CHANNEL}'", "user_id": "'${SLACK_USERID}'"}}' \
  https://${github_base}/repos/${REPO_OWNER}/${REPO_NAME}/disapatches 
