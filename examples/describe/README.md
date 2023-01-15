# Bashbot Example - Describe Command

In this example, the configuration yaml file for Bashbot is parsed via [jq](https://stedolan.github.io/jq/) to display the trigger and name of each command in the file.

<img src="https://i.imgur.com/HAJ3TS2.gif">

## Bashbot Configuration

This command is triggered by sending `!bashbot describe [command]` in a slack channel where Bashbot is also a member. There is no external script for this command, and expects yq to already be installed, and the environment variable `BASHBOT_CONFIG_FILEPATH` to be pointing at the running configuration, on the host machine.

```yaml
name: Describe Bashbot [command]
description: Show the yaml object for a specific command
envvars:
  - BASHBOT_CONFIG_FILEPATH
dependencies:
  - yq
help: "!bashbot describe [command]"
trigger: describe
location: /bashbot/
command:
  - yq e '.tools[] | select(.trigger=="${command}")' ${BASHBOT_CONFIG_FILEPATH}
parameters:
  - name: command
    allowed: []
    description: a command to describe ('bashbot list')
    source:
      - yq e '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}
log: true
ephemeral: false
response: code
permissions:
  - all
```
