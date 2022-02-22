# Example Bashbot Commands

This directory contains a few sample Bashbot commands that can be applied to a configuration json by executing a helper script or by copy/pasting the json directly. Each example (directory) also contains a read-me file describing usage information. The [add-example/add-command.sh](add-example/add-command.sh) script expects the environment variable `BASH_CONFIG_FILEPATH` to be set and takes one argument (filepath of the command to add to the configuration json). The [remove-example/remove-command.sh](remove-example/remove-command.sh) script also expects the environment variable `BASH_CONFIG_FILEPATH` to be set and takes one argument (the value of the parameter `trigger` within the specific command).


```bash
# From Bashbot source root
export BASHBOT_CONFIG_FILEPATH=${PWD}/config.json

./add-example/add-command.sh version/version.json
# version added to /Users/user/bashbot/config.json

./add-example/add-command.sh version/version.json
# Trigger already exists: version

./remove-example/remove-command.sh version
# version removed from /Users/user/bashbot/config.json

./remove-example/remove-command.sh version
# Trigger not found to remove: version
```

## Examples


1. Simple call/response
    - [Ping/Pong](ping)
        - \+ `./add-example/add-command.sh ping/ping.json`
        - \- `./remove-example/remove-command.sh ping`
2. Run bash script
    - [Get Caller Information](info)
        - \+ `./add-example/add-command.sh info/info.json`
        - \- `./remove-example/remove-command.sh info`
    - [Get asdf Plugins](asdf)
        - \+ `./add-example/add-command.sh asdf/asdf.json`
        - \- `./remove-example/remove-command.sh asdf`
    - [Get add-example Plugins](add-example)
        - \+ `./add-example/add-command.sh add-example/add-example.json`
        - \- `./remove-example/remove-command.sh add-example`
    - [Get remove-example Plugins](remove-example)
        - \+ `./add-example/add-command.sh remove-example/remove-example.json`
        - \- `./remove-example/remove-command.sh remove-example`
    - [Get list-examples Plugins](list-examples)
        - \+ `./add-example/add-command.sh list-examples/list-examples.json`
        - \- `./remove-example/remove-command.sh list-examples`
3. Run go-lang script
    - [Get Running Version](version)
        - \+ `./add-example/add-command.sh version/version.json`
        - \- `./remove-example/remove-command.sh version`
4. Parse json via jq
    - [Show Help Dialog](help)
        - \+ `./add-example/add-command.sh help/help.json`
        - \- `./remove-example/remove-command.sh help`
    - [List Available Commands](list)
        - \+ `./add-example/add-command.sh list/list.json`
        - \- `./remove-example/remove-command.sh list`
    - [Describe Command](describe)
        - \+ `./add-example/add-command.sh describe/describe.json`
        - \- `./remove-example/remove-command.sh describe`
5. Curl/wget
    - [Get Latest Bashbot Release](latest-release)
        - \+ `./add-example/add-command.sh latest-release/latest-release.json`
        - \- `./remove-example/remove-command.sh latest-release`
    - [Get File From Repository](get-file-from-repo)
        - \+ `./add-example/add-command.sh get-file-from-repo/get-file-from-repo.json`
        - \- `./remove-example/remove-command.sh get-file-from-repo`
    - [Trigger Github Action](trigger-github-action)
        - \+ `./add-example/add-command.sh trigger-github-action/trigger-github-action.json`
        - \- `./remove-example/remove-command.sh trigger-github-action`
    - [Get Air Quality Index By Zip](aqi)
        - \+ `./add-example/add-command.sh aqi/aqi.json`
        - \- `./remove-example/remove-command.sh aqi`
6. [Send message from github action](#send-message-from-github-action)


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
      "install": [
        "rm -rf bashbot-scripts || true",
        "git clone https://github.com/mathew-fleisch/bashbot-scripts.git"
      ]
    }
  ]
}
```

Each object in the tools array defines the parameters of a single command.

```
name, description and help provide human readable information about the specific command
trigger:      unique alphanumeric word that represents the command
location:     absolute or relative path to dependency directory (use "./" for no dependency)
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


-----------------------


## Send message from github action

Bashbot's binary can also be used to send a single message to a slack channel. This can be useful at the end of existing automation as a notification of success/failure. In this example, a [github-action](../.github/workflows/example-notify-slack.yaml) is triggered via curl. The [github-action itself](../.github/workflows/example-notify-slack.yaml) will [install bashbot via asdf](https://github.com/mathew-fleisch/asdf-bashbot/), then send the message from the curl's payload to a specific channel also defined in the payload.

```bash
#!/bin/bash
# shellcheck disable=SC2086

REPO_OWNER=mathew-fleisch
REPO_NAME=bashbot
SLACK_CHANNEL=GPFMM5MD2
SLACK_MESSAGE="Custom Notification from Github Action"

github_base="${github_base:-api.github.com}"
expected_variables="SLACK_CHANNEL SLACK_MESSAGE REPO_OWNER REPO_NAME GITHUB_TOKEN"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 0
  fi
done

DATA='{"event_type":"trigger-slack-notify","client_payload": {"channel":"'${SLACK_CHANNEL}'", "text":"'${SLACK_MESSAGE}'"}}'
DATA=$(echo "${DATA}" | jq -c .)
curl \
  -X POST \
  -H "Accept: application/vnd.github.everest-preview+json" \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  --data "${DATA}" \
  https://${github_base}/repos/${REPO_OWNER}/${REPO_NAME}/dispatches

```

The above script will trigger the [github action](../.github/workflows/example-notify-slack.yaml) (also copied below). After installing bashbot via asdf, bashbot's binary can be executed using environment variables and passing values from the curl's payload. This method of installing executing Bashbot can be used in any github action pipeline.

```yaml
# Name:        example-notify-slack.yaml
# Author:      Mathew Fleisch <mathew.fleisch@gmail.com>
# Description: This action demonstrates how to trigger a Slack Notification from Bashbot.
name: Example Bashbot Notify Slack
on:
  repository_dispatch:
    types:
      - trigger-slack-notify

jobs:
  build:
    name: Example Bashbot Notify Slack
    runs-on: ubuntu-latest
    steps:
      -
        name: Install Bashbot via asdf
        uses: asdf-vm/actions/install@v1
        with:
          tool_versions: bashbot 1.8.0
      -
        name: Send Slack Message With Bashbot Binary
        env:
          BASHBOT_CONFIG_FILEPATH: ./config.json
          SLACK_TOKEN: ${{ secrets.RELEASE_SLACK_TOKEN }}
        run: |
          echo '{"admins":[{"trigger":"bashbotexample","appName":"Bashbot Example","userIds":[""],"privateChannelId":"","logChannelId":""}],"messages":[],"tools":[],"dependencies":[]}' > $BASHBOT_CONFIG_FILEPATH
          bashbot \
            --send-message-channel ${{ github.event.client_payload.channel }} \
            --send-message-text "${{ github.event.client_payload.text }}"
```

Executing this action has logs/output that looks similar to this

<img src="https://i.imgur.com/blZLZmR.png" />

<img src="https://i.imgur.com/kkaClWJ.png" />