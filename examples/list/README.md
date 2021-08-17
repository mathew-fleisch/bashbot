# Bashbot Example - List Commands

In this example, the configuration json file for Bashbot is parsed via [jq](https://stedolan.github.io/jq/) to display the trigger and name of each command in the file.

<img src="https://i.imgur.com/HHzHlFK.gif">

## Bashbot Configuration

This command is triggered by sending `bashbot list` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects jq to already be installed on the host machine.

```json
{
  "name": "List Available Bashbot Commands",
  "description": "List all of the possible commands stored in bashbot",
  "help": "bashbot list",
  "trigger": "list",
  "location": "./",
  "command": ["jq -r '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}"],
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```