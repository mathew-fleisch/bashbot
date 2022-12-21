# Bashbot Example - Trigger Github Action

In this example, a curl is executed via bash script and triggers a github action. The github action downloads the latest version of the bashbot binary to reply back to the caller upon completion of the job. This command has four parts: Bashbot configuration, curl trigger, github-action yaml, github-action script

<img src="https://i.imgur.com/s0cf2Hl.gif" />

## Bashbot configuration

This command is triggered by sending `bashbot trigger-github-action` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/trigger-github-action/trigger-github-action` and requires the following environment variables to be set: `BASHBOT_CONFIG_FILEPATH SLACK_TOKEN SLACK_CHANNEL SLACK_USERID REPO_OWNER REPO_NAME GITHUB_TOKEN GITHUB_RUN_ID`.

```json
{
  "name": "Trigger a Github Action",
  "description": "Triggers an example Github Action job by repository dispatch",
  "envvars": [
    "BASHBOT_CONFIG_FILEPATH", 
    "SLACK_TOKEN", 
    "SLACK_CHANNEL", 
    "SLACK_USERID", 
    "REPO_OWNER",
    "REPO_NAME",
    "GITHUB_TOKEN",
    "GITHUB_RUN_ID"
  ],
  "dependencies": ["curl", "wget"],
  "help": "bashbot trigger-github-action",
  "trigger": "trigger-github-action",
  "location": "./examples/github-action",
  "command": [
    "export REPO_OWNER=mathew-fleisch",
    "&& export REPO_NAME=bashbot",
    "&& export SLACK_CHANNEL=${TRIGGERED_CHANNEL_ID}",
    "&& export SLACK_USERID=${TRIGGERED_USER_ID}",
    "&& ./trigger.sh",
    "&& echo \"Running this <https://github.com/${REPO_OWNER}/${REPO_NAME}/blob/main/.github/workflows/example-bashbot-github-action.yaml|example github action>\""
  ],
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "text",
  "permissions": ["all"]
}
```

## Bashbot scripts

There are two scripts associated with this Bashbot command: [trigger.sh](trigger.sh) and [github-action.sh])(github-action.sh). The [trigger.sh](trigger.sh) script sends off a POST request via curl to the repository's dispatch function to trigger a github action. The [github-action](../.github/workflows/example-bashbot-github-action.yaml) uses the [github-action.sh])(github-action.sh) script to simulate a long running job, and return back status to slack via Bashbot binary.

### trigger.sh

This script is used to trigger the github action.

```bash
github_base="${github_base:-api.github.com}"
expected_variables="BASHBOT_CONFIG_FILEPATH SLACK_TOKEN SLACK_CHANNEL SLACK_USERID REPO_OWNER REPO_NAME GITHUB_TOKEN"
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
  -H "Authorization: token ${GITHUB_TOKEN}" \
  --data '{"event_type":"trigger-github-action","client_payload": {"channel":"'${SLACK_CHANNEL}'", "user_id": "'${SLACK_USERID}'"}}' \
  "https://${github_base}/repos/${REPO_OWNER}/${REPO_NAME}/dispatches"

```


### github-action.sh

This script uses curl to get the latest version and job id from the github api, then uses wget to get the latest version of bashbot. Finally this script uses the Bashbot binary to send a message back to slack (simulating a long running job).

```bash
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
JOB_ID=$(curl ${headers} https://${github_base}/repos/${REPO_OWNER}/${REPO_NAME}/actions/runs/${GITHUB_RUN_ID}/jobs | jq -r '.jobs[].id')
arch=amd64
[ $(uname -m) == "aarch64" ] && arch=arm64
os=$(uname | tr '[:upper:]' '[:lower:]')
wget -q https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/bashbot-${os}-${arch} -O bashbot
chmod +x bashbot
downloaded_version=$(./bashbot version | awk '{print $2}')
./bashbot send-message \
    --channel ${SLACK_CHANNEL} \
    --msg "<@${SLACK_USERID}> Bashbot triggered this <https://github.com/${REPO_OWNER}/${REPO_NAME}/blob/main/.github/workflows/example-bashbot-github-action.yaml|example github action> and used the bashbot binary (<https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/tag/${downloaded_version}|${downloaded_version}>) within the github action, to <https://github.com/${REPO_OWNER}/${REPO_NAME}/runs/${JOB_ID}?check_suite_focus=true|simulate a long running job> in order to send success/failure back to slack."
```