# Bashbot Example - List Commands

In this example, the configuration yaml file for Bashbot is parsed via yq to display the trigger and name of each command in the file.

<img src="https://i.imgur.com/HHzHlFK.gif">

## Bashbot Configuration

This command is triggered by sending `!bashbot list` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects yq to already be installed, and the environment variable `BASHBOT_CONFIG_FILEPATH` to be pointing at the running configuration, on the host machine.

```yaml
name: List Available Bashbot Commands
description: List all of the possible commands stored in bashbot
envvars:
  - BASHBOT_CONFIG_FILEPATH
dependencies:
  - yq
help: "!bashbot list"
trigger: list
location: /bashbot/
command:
  - yq e '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}
parameters: []
log: true
ephemeral: false
response: code
permissions:
  - all
```
