# Bashbot Example - List Example Commands

In this example, all of the yaml filenames are aggregated and sorted as a list.

<img src="https://i.imgur.com/e4nEuSE.gif">

## Bashbot Configuration

This command is triggered by sending `!bashbot list-examples` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects Bashbot's examples directory to exist. This command takes no arguments.

```yaml
name: List Example Commands
description: List commands from bashbot example commands
envvars: []
dependencies: []
help: "!bashbot list-examples"
trigger: list-examples
location: /bashbot/vendor/bashbot/examples
command:
  - find . -name "*.json"
  - "| xargs -I {} bash -c"
  - "'export example=$(basename {} .json)"
  - "&& printf \"%21s - %s\" \"$example\" \"https://github.com/mathew-fleisch/bashbot/tree/main/examples/$example\""
  - "&& echo'"
  - "| sort -k 2"
parameters: []
log: true
ephemeral: false
response: code
permissions:
  - all
```
