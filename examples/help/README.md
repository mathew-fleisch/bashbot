# Bashbot Example - Help Dialog

In this example, the configuration yaml file for Bashbot is parsed via yq to display the help and description values for each command

<img src="https://i.imgur.com/QC2av32.gif" />

## Bashbot Configuration

This command is triggered by sending `!bashbot help` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects yq to already be installed, and the environment variable `BASHBOT_CONFIG_FILEPATH` to be pointing at the running configuration, on the host machine.

```yaml
name: BashBot Help
description: Show this message
envvars:
  - BASHBOT_CONFIG_FILEPATH
dependencies:
  - yq
help: "!bashbot help"
trigger: help
location: /bashbot/
command:
  - echo "BashBot is a tool for infrastructure/devops teams to automate tasks triggered by slash-command-like declarative configuration" &&
  - echo '```' &&
  - "yq e '.tools[] | {.help: .description}' \"${BASHBOT_CONFIG_FILEPATH}\""
  - "| sed -e 's/\\\"//g'"
  - "| sed -e 's/:/ -/g' &&"
  - echo '```'
parameters: []
log: true
ephemeral: false
response: text
permissions:
  - all
```
