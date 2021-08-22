# Bashbot Example - Add Example Command

In this example, a json example can be added to a running configuration json. 

***Note: This will not work when the configuration json is mounted as a configmap. Use the seed method if bashbot is deployed in kubernetes to use this example***

## Bashbot Configuration

This command is triggered by sending `bashbot asdf` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects Bashbot's examples directory to exist. The valid arguments for this command come from the output of `bashbot list-examples`

```json
{
  "name": "Add Example Command",
  "description": "Add command from Bashbot example commands",
  "help": "bashbot add-example",
  "trigger": "add-example",
  "location": "./examples",
  "command": ["./add-example/add-command.sh $(find . -name \"${add_command}.json\")"],
  "parameters": [
    {
      "name": "add_command",
      "allowed": [],
      "description": "a command to add to bashbot config ('bashbot list-examples')",
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