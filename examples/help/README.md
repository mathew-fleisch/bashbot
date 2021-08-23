# Bashbot Example - Help Dialog

In this example, the configuration json file for Bashbot is parsed via [jq](https://stedolan.github.io/jq/) to display the help and description values for each command

## Bashbot Configuration

This command is triggered by sending `bashbot help` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects jq to already be installed, and the environment variable `BASHBOT_CONFIG_FILEPATH` to be pointing at the running configuration, on the host machine.

```json
{
  "name": "BashBot Help",
  "description": "Show this message",
  "help": "bashbot help",
  "trigger": "help",
  "location": "./",
  "command": ["jq -r '.tools[] | \"\\(.help) - \\(.description)\"' ${BASHBOT_CONFIG_FILEPATH}"],
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": [
    "all"
  ]
}
```