# Bashbot Example - Latest Release

In this example, a curl is used to determine the latest version of Bashbot through the github release api

<img src="https://i.imgur.com/w3wouOR.gif">

## Bashbot configuration

This command is triggered by sending `!bashbot latest-release` in a slack channel where Bashbot is also a member. This command requires no external script and simply returns a formatted message including links to Bashbot's source code and latest release. This command requires [curl](https://curl.se/) to be installed on the host machine.

```yaml
name: Get Latest Bashbot Version
description: Returns the latest version of Bashbot via curl
envvars: []
dependencies:
  - curl
help: "!bashbot latest-release"
trigger: latest-release
location: /bashbot/
command:
  - latest_version=$(curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)
  - "&& echo \"The latest version of <https://github.com/mathew-fleisch/bashbot|Bashbot>: <https://github.com/mathew-fleisch/bashbot/releases/tag/$latest_version|$latest_version>\""
parameters: []
log: true
ephemeral: false
response: text
permissions:
  - all
```
