# Bashbot Example - Trigger Github Action

In this example, a curl is executed via bash script and triggers a github action. The github action downloads the latest version of the bashbot binary to reply back to the caller upon completion of the job. This command has four parts: Bashbot configuration, curl trigger, github-action yaml, github-action script.

<img src="https://i.imgur.com/s0cf2Hl.gif" />

## Bashbot configuration

This command is triggered by sending `!bashbot trigger-github-action` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/trigger-github-action` and requires the following environment variables to be set: `BASHBOT_CONFIG_FILEPATH SLACK_BOT_TOKEN SLACK_APP_TOKEN SLACK_CHANNEL SLACK_USERID REPO_OWNER REPO_NAME GITHUB_TOKEN GITHUB_RUN_ID`. The permissions array restricts this command to only run in a specific slack channel.

```yaml
name: Trigger a Github Action
description: Triggers an example Github Action job by repository dispatch
envvars:
  - GIT_TOKEN
dependencies: []
help: "!bashbot trigger-github-action"
trigger: trigger-github-action
location: /bashbot/vendor/bashbot/examples/trigger-github-action
command:
  - "export REPO_OWNER=mathew-fleisch"
  - "&& export REPO_NAME=bashbot"
  - "&& export SLACK_CHANNEL=${TRIGGERED_CHANNEL_ID}"
  - "&& export SLACK_USERID=${TRIGGERED_USER_ID}"
  - "&& echo \"Running this example github action: https://github.com/${REPO_OWNER}/${REPO_NAME}/blob/main/.github/workflows/example-bashbot-github-action.yaml\""
  - "&& ./trigger.sh"
parameters: []
log: true
ephemeral: false
response: text
permissions:
  - GPFMM5MD2
```

## Bashbot scripts

There are two scripts associated with this Bashbot command: [trigger.sh](trigger.sh) and [github-action.sh](github-action.sh). The [trigger.sh](trigger.sh) script sends off a POST request via curl to the repository's dispatch function to trigger a github action. The [github-action](../.github/workflows/example-bashbot-github-action.yaml) uses the [github-action.sh](github-action.sh) script to simulate a long running job, and return back status to slack via Bashbot binary.

### trigger.sh and trigger-gate.sh

The scripts [examples/trigger-github-action/trigger.sh](examples/trigger-github-action/trigger.sh) and [examples/trigger-github-action/trigger-gate.sh](examples/trigger-github-action/trigger-gate.sh), use curl to trigger one of two github actions. [.github/workflows/example-bashbot-github-action.yaml](.github/workflows/example-bashbot-github-action.yaml) shows how bashbot can be used to trigger a github action, and also used to craft a custom pass/failure message back to slack.  [.github/workflows/example-bashbot-github-action-gate.yaml](.github/workflows/example-bashbot-github-action-gate.yaml) shows how bashbot can be used to trigger a github action, and also used as a gating mechanism to pass/fail the job manually from slack.
