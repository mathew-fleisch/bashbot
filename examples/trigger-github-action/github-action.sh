#!/bin/bash

github_base="${github_base:-api.github.com}"
expected_variables="BASHBOT_CONFIG_FILEPATH SLACK_TOKEN SLACK_CHANNEL SLACK_USERID REPO_OWNER REPO_NAME GITHUB_TOKEN GITHUB_RUN_ID"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 1
  fi
done
headers="-sH \"Accept: application/vnd.github.everest-preview+json\" -H \"Authorization: token ${GITHUB_TOKEN}\""
LATEST_VERSION=$(curl ${headers} https://${github_base}/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest | grep tag_name | cut -d '"' -f 4)
arch=amd64
[ $(uname -m) == "aarch64" ] && arch=arm64
os=$(uname | tr '[:upper:]' '[:lower:]')
wget -q https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/bashbot-${os}-${arch} -O bashbot
chmod +x bashbot
./bashbot \
    --send-message-channel ${SLACK_CHANNEL} \
    --send-message-text "<@${SLACK_USERID}> Bashbot triggered this job: \${{ github.server_url }}/\${{ github.repository }}/actions/runs/\${{ github.run_id }}"