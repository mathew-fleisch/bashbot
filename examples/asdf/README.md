# Bashbot Example - List asdf Plugins

In this example, a few asdf commands are combined to display the installed versions for each of the pre-installed ([in container](.tool-versions)) asdf plugins.

## Bashbot Configuration

This command is triggered by sending `bashbot asdf` in a slack channel where Bashbot is also a member. There is no external script for this command, takes no arugments/parameters, and expects asdf and asdf plugins to already be installed on the host machine. The command first sources asdf, and echos some title headers. Then it pipes the output of `asdf plugin list` to `xargs` in order to get the installed version for each plugin. It then uses `printf` to print the asdf plugin name and its version, in a table format. Escaping single and double quotes in bashbot's json syntax can be difficult.

```json
{
  "name": "List asdf plugins",
  "description": "List the pre-installed asdf plugins",
  "help": "bashbot asdf",
  "trigger": "asdf",
  "location": "./",
  "command": [
    ". $ASDF_DATA_DIR/asdf.sh",
    "&& echo \"asdf plugin | version\"",
    "&& echo \"------------|--------\"",
    "&& asdf plugin list | xargs -I {} bash -c 'printf \"%10s  | %s\" \"{}\" \"$(asdf list {} | awk '\"'\"'{print $1}'\"'\"')\" && echo'"
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