# Bashbot Example - List asdf Plugins

In this example, a few asdf commands are combined to display the installed versions for each of the pre-installed ([in container](.tool-versions)) asdf plugins.

<img src="https://i.imgur.com/v1aqdj6.png" />

## Bashbot Configuration

This command is triggered by sending `bashbot asdf` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects asdf and asdf plugins to already be installed on the host machine. The command first sources asdf, and echos some title headers. Then it pipes the output of `asdf plugin list` to `xargs` in order to get the installed version for each plugin. It then uses `printf` (along with `awk` to remove leading spaces from `asdf list [ASDF-PLUGIN-NAME]` command) to print each asdf plugin name and its version, in a table format. Escaping single and double quotes in bashbot's json syntax can be difficult.

```json
{
  "name": "List asdf plugins",
  "description": "List the installed asdf plugins and their versions",
  "help": "bashbot asdf",
  "trigger": "asdf",
  "location": "./",
  "command": [
    ". $ASDF_DATA_DIR/asdf.sh",
    "&& echo \"•──────────────────────────────•\"",
    "&& echo \"│ <https://asdf-vm.com/|asdf version: $(asdf version)> |\"",
    "&& echo \"├───────────────────┰──────────┤\"",
    "&& echo \"│       asdf plugin │ version  │\"",
    "&& echo \"├───────────────────┼──────────┤\"",
    "&& asdf plugin list",
      "| xargs -I {} bash -c",
        "'printf \"│ %17s │ %-8s │\"",
        "\"{}\"",
        "\"$(asdf list {} | awk '\"'\"'{print $1}'\"'\"')\"",
        "&& echo'",
    "&& echo \"•──────────────────────────────•\""
  ],
  "parameters": [],
  "log": true,
  "ephemeral": false,
  "response": "code",
  "permissions": [
    "all"
  ]
}
```