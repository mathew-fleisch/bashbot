# Bashbot Example - List asdf Plugins

In this example, a few asdf commands are combined to display the installed versions for each of the pre-installed ([in container](../../.tool-versions)) asdf plugins.

<img src="https://i.imgur.com/v1aqdj6.png" />

## Bashbot Configuration

This command is triggered by sending `!bashbot asdf` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects asdf and asdf plugins to already be installed on the host machine/container.

```yaml
name: List asdf plugins
description: List the installed asdf plugins and their versions
envvars: []
dependencies: []
help: "!bashbot asdf"
trigger: asdf
location: /bashbot/
command:
  - ". /usr/asdf/asdf.sh && asdf list"
parameters: []
log: true
ephemeral: false
response: code
permissions:
  - all
```
