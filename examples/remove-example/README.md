# Bashbot Example - Remove Example Command

In this example, a json example can be removed from a running configuration. 

***Note: This will not work when the configuration json is mounted as a configmap. Use the seed method if bashbot is deployed in kubernetes to use this example***

## Bashbot Configuration

This command is triggered by sending `bashbot remove-example [command]` in a slack channel where Bashbot is also a member. There is no external script for this command and expects Bashbot's examples directory to exist. The valid arguments for this command come from the output of `bashbot list-examples` [list-examples read-me](../list-examples)

```json
{
  "name": "Remove Example Command",
  "description": "Remove command from bashbot example commands",
  "help": "bashbot remove-example",
  "trigger": "remove-example",
  "location": "./examples",
  "command": ["./remove-command.sh ${remove_command}"
  ],
  "parameters": [
    {
      "name": "remove_command",
      "allowed": [],
      "description": "a command to remove to bashbot config ('bashbot list-examples')",
      "source": ["find . -name \"*.json\" | xargs -I {} basename {} .json"]
    }
  ],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": [
    "all"
  ]
}
```