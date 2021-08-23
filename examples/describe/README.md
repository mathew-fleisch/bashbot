# Bashbot Example - Describe Command

In this example, the configuration json file for Bashbot is parsed via [jq](https://stedolan.github.io/jq/) to display the trigger and name of each command in the file.

<img src="https://i.imgur.com/bQZKRjX.gif">

## Bashbot Configuration

This command is triggered by sending `bashbot describe [command]` in a slack channel where Bashbot is also a member. There is no external script for this command, and expects jq to already be installed, and the environment variable `BASHBOT_CONFIG_FILEPATH` to be pointing at the running configuration, on the host machine. One argument/parameter is required to run this command and is used to select and print specific commands from the configuration json using jq. The `command` parameter is used to  `.tools[].trigger`

```json
{
  "name": "Describe Bashbot [command]",
  "description": "Show the json object for a specific command",
  "help": "bashbot describe [command]",
  "trigger": "describe",
  "location": "./",
  "command": ["jq '.tools[] | select(.trigger==\"${command}\")' ${BASHBOT_CONFIG_FILEPATH}"],
  "parameters": [
    {
      "name": "command",
      "allowed": [],
      "description": "a command to describe ('bashbot list')",
      "source": ["jq -r '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}"]
    }
  ],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```