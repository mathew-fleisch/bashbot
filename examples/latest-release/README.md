# Bashbot Example - Latest Release

In this example, a curl is used to determine the latest version of Bashbot through the github release api

## Bashbot configuration

This command is triggered by sending `bashbot latest-release` in a slack channel where Bashbot is also a member. This command requires no external script and simply returns a formatted message including links to Bashbot's source code and latest release.

```json
{
  "name": "Get Latest Bashbot Version",
  "description": "Returns the latest version of Bashbot via curl",
  "help": "bashbot latest-release",
  "trigger": "latest-release",
  "location": "./",
  "setup": "latest_version=$(curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '\"' -f 4)",
  "command": "echo \"The latest version of <https://github.com/mathew-fleisch/bashbot|Bashbot>: <https://github.com/mathew-fleisch/bashbot/releases/tag/$latest_version|$latest_version>\"",
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "text",
  "permissions": ["all"]
}
```

