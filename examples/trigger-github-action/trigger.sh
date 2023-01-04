#!/bin/bash
# shellcheck disable=SC2086

github_base="${github_base:-api.github.com}"
expected_variables="BASHBOT_CONFIG_FILEPATH SLACK_CHANNEL SLACK_USERID REPO_OWNER REPO_NAME GIT_TOKEN"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 0
  fi
done

curl -s \
  -X POST \
  -H "Accept: application/vnd.github.everest-preview+json" \
  -H "Authorization: token ${GIT_TOKEN}" \
  --data '{"event_type":"trigger-github-action","client_payload": {"channel":"'${SLACK_CHANNEL}'", "user_id": "'${SLACK_USERID}'"}}' \
  "https://${github_base}/repos/${REPO_OWNER}/${REPO_NAME}/dispatches"
