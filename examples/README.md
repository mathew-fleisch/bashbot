# Example Bashbot Commands

This directory contains a few sample Bashbot commands that can be applied to a configuration json by executing a helper script or by copy/pasting the json directly. Each example (directory) also contains a read-me file describing usage information. The [add-command.sh](add-command.sh) script expects the environment variable `BASH_CONFIG_FILEPATH` to be set and takes one argument (filepath of the command to add to the configuration json). The [remove-command.sh](remove-command.sh) script also expects the environment variable `BASH_CONFIG_FILEPATH` to be set and takes one argument (the value of the parameter `trigger` within the specific command).


```bash
# From Bashbot source root
export BASHBOT_CONFIG_FILEPATH=${PWD}/config.json

./examples/add-command.sh examples/version/version.json              
# version added to /Users/user/bashbot/config.json

./examples/add-command.sh examples/version/version.json                                           
# Trigger already exists: version

./examples/remove-command.sh version                
# version removed from /Users/user/bashbot/config.json

./examples/remove-command.sh version
# Trigger not found to remove: version
```

## Examples


1. Simple call/response
    - [Ping/Pong](ping)
        - \+ `./add-command.sh ping/ping.json`
        - \- `./remove-command.sh ping`
2. Run bash script
    - [Get Caller Information](info)
        - \+ `./add-command.sh info/info.json`
        - \- `./remove-command.sh info`
3. Run golang script
    - [Get Running Version](version)
        - \+ `./add-command.sh version/version.json`
        - \- `./remove-command.sh version`
4. Parse json via jq
    - [List Available Commands](list)
        - \+ `./add-command.sh list/list.json`
        - \- `./remove-command.sh list`
    - [Describe Command](describe)
        - \+ `./add-command.sh describe/describe.json`
        - \- `./remove-command.sh describe`
5. Curl/wget
    - [Get Latest Bashbot Release](latest-release)
        - \+ `./add-command.sh latest-release/latest-release.json`
        - \- `./remove-command.sh latest-release`
    - [Get File From Repository](get-file-from-repo)
        - \+ `./add-command.sh get-file-from-repo/get-config.json`
        - \- `./remove-command.sh get-config`
    - [Trigger Github Action](github-action)
        - \+ `./add-command.sh github-action/trigger-github-action.json`
        - \- `./remove-command.sh trigger-github-action`



-------------------------------------------------------------------------

### config.json
[sample-config.json](../sample-config.json)
The config.json file is defined as an array of json objects keyed by 'tools' and 'dependencies.' The dependencies section defines any resources that need to be downloaded or cloned from a repository before execution of commands. The following is a simplified example of a config.json file:

```json
{
  "admins": [{
    "trigger": "bashbot",
    "appName": "BashBot",
    "userIds": ["SLACK-USER-ID"],
    "privateChannelId": "SLACK-CHANNEL-ID",
    "logChannelId": "SLACK-CHANNEL-ID"
  }],
  "messages": [{
    "active": true,
    "name": "welcome",
    "text": "Witness the power of %s"
  },{
    "active": true,
    "name": "processing_command",
    "text": ":robot_face: Processing command..."
  },{
    "active": true,
    "name": "processing_raw_command",
    "text": ":smiling_imp: Processing raw command..."
  },{
    "active": true,
    "name": "command_not_found",
    "text": ":thinking_face: Command not found..."
  },{
    "active": true,
    "name": "incorrect_parameters",
    "text": ":face_with_monocle: Incorrect number of parameters"
  },{
    "active": true,
    "name": "invalid_parameter",
    "text": ":face_with_monocle: Invalid parameter value: %s"
  },{
    "active": true,
    "name": "ephemeral",
    "text": ":shushing_face: Message only shown to user who triggered it."
  },{
    "active": true,
    "name": "unauthorized",
    "text": ":skull_and_crossbones: You are not authorized to use this command in this channel.\nAllowed in: [%s]"
  }],
  "tools": [{
      "name": "List Commands",
      "description": "List all of the possible commands stored in bashbot",
      "help": "bashbot list-commands",
      "trigger": "list-commands",
      "location": "./",
      "setup": "echo \"\"",
      "command": ["jq -r '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}"],
      "parameters": [],
      "log": false,
      "ephemeral": false,
      "response": "code",
      "permissions": ["all"]
    }
  ],
  "dependencies": [
    {
      "name": "BashBot scripts Scripts",
      "install": "rm -rf bashbot-scripts || true && git clone https://github.com/mathew-fleisch/bashbot-scripts.git"
    }
  ]
}
```

Each object in the tools array defines the parameters of a single command.

```
name, description and help provide human readable information about the specific command
trigger:      unique alphanumeric word that represents the command
location:     absolute or relative path to dependency directory (use "./" for no dependency)
setup:        command that is run before the main command. (use "echo \"\"" as a default)
command:      bash command using ${parameter-name} to inject white-listed parameters or environment variables
parameters:   array of parameter objects. (more detail below)
log:          define whether the command should be logged in log channel
ephemeral:    define if the response should be shown to all, or just the user that triggered the command
response:     [code|text] code displays response in a code block, text displays response as raw text
permissions:  array of strings. private channel ids to restrict command access to
```

#### parameters
Parameters passed to Bashbot cannot be arbitrary/free-form text and must come from a hard-coded list of values, or a dynamic list of values as the return of another command. In this first example a hard-coded list of values is used to print the current `date` or `uptime` by passing to bashbot `bashbot time date` or `bashbot time uptime`

```json
{
  "name": "Date or Uptime",
  "description": "Show the current time or uptime",
  "help": "bashbot time",
  "trigger": "time",
  "location": "./",
  "setup": "echo \"\"",
  "command": ["${command}"],
  "parameters": [
    {
      "name": "command",
      "allowed": ["date","uptime"]
    }
  ],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```

In this next example, the command is derived by the output of another command. The valid parameters in this case, are the "trigger" values at `.tools[].trigger` and are used to output the json object, through jq, at that array index. Running `bashbot describe describe` would therefore output itself.

```json
{
  "name": "Describe Bashbot [command]",
  "description": "Show the json object for a specific command",
  "help": "bashbot describe [command]",
  "trigger": "describe",
  "location": "./",
  "setup": "echo \"\"",
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