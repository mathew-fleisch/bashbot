# Bashbot Example - List Example Commands

In this example, all of the json filenames are aggregated and sorted as a list. The values listed can be passed to the commands [add-example](../add-example) and [remove-example](../remove-example)

***Note: This will not work when the configuration json is mounted as a configmap. Use the seed method if bashbot is deployed in kubernetes to use this example***

## Bashbot Configuration

This command is triggered by sending `bashbot list-examples` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects Bashbot's examples directory to exist. This command takes no arguments.

```json
{
  "name": "List Example Commands",
  "description": "List commands from bashbot example commands",
  "envvars": [],
  "dependencies": [],
  "help": "bashbot list-examples",
  "trigger": "list-examples",
  "location": "./examples",
  "command": [
    "find . -name \"*.json\"",
    "| xargs -I {} bash -c",
      "'export example=$(basename {} .json)",
      "&& printf \"%21s - %s\" \"$example\" \"https://github.com/mathew-fleisch/bashbot/tree/main/examples/$example\"",
      "&& echo'",
    "| sort -k 2"
  ],
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": [
    "all"
  ]
}
```
