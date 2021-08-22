# Bashbot Example - Add Example Command

In this example, a json example can be added to a running configuration json. 

***Note: This will not work when the configuration json is mounted as a configmap. Use the seed method if bashbot is deployed in kubernetes to use this example***

## Bashbot Configuration

This command is triggered by sending `bashbot asdf` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects Bashbot's examples directory to exist. This command takes no arguments.

```json
{
  "name": "List Example Commands",
  "description": "List commands from bashbot example commands",
  "help": "bashbot list-examples",
  "trigger": "list-examples",
  "location": "./examples",
  "command": [
    "find . -name \"*.json\" | xargs -I {} basename {} .json"
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