# Bashbot

[![Build binaries](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml)
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml)
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated)

BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to determine when to trigger bash scripts. A [configuration file](sample-config.json) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls.


--------------------------------------------------

## Installation and setup 

Bashbot can be run as a go binary or as a container and requires a slack-token and a config.json. The go binary takes flags to set the slack-token and path to the config.json file and the container uses environment variables to trigger a go binary by [entrypoint.sh](entrypoint.sh). 

***Note about slack-token***

Slack's permissions model for the "[Real-Time-Messaging (RTM)](https://api.slack.com/rtm)" socket connection, requires a "classic app" to be configured to get the correct type of token to run Bashbot. After logging into slack via browser, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "legacy bot user" and "Bot User OAuth Access Token." Finally, add bashbot to your workspace and invite to a channel.


### Run Bashbot as a go-binary

***Requirements***

- jq
- git
- golang
- wget
- curl

```bash
# Download and install the latest version of bashbot
make install-latest
# or manually (without golang)
os=$(uname | tr '[:upper:]' '[:lower:]')
arch="amd64"
if [ "$(uname -m)" == "aarch64" ]; then arch="arm64"; fi
latest=$(curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)
wget -O /usr/local/bin/bashbot https://github.com/mathew-fleisch/bashbot/releases/download/${latest}/bashbot-${os}-${arch}
chmod +x /usr/local/bin/bashbot

# Set environment variables
export SLACK_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
export BASHBOT_CONFIG_FILEPATH=${PWD}/config.json

# Get the version
bashbot --version
# bashbot-darwin-amd64  v1.4.6

# Show the help dialog
bashbot --help
#  ____            _     ____        _   
# |  _ \          | |   |  _ \      | |  
# | |_) | __ _ ___| |__ | |_) | ___ | |_ 
# |  _ < / _' / __| '_ \|  _ < / _ \| __|
# | |_) | (_| \__ \ | | | |_) | (_) | |_ 
# |____/ \__,_|___/_| |_|____/ \___/ \__|
# Bashbot is a slack bot, written in golang, that can be configured
# to run bash commands or scripts based on a configuration file.

# Usage: ./bashbot [flags]

#   -config-file string
#         [REQUIRED] Filepath to config.json file
#   -help
#         Help/usage information
#   -log-format string
#         Display logs as json or text (default "text")
#   -log-level string
#         Log level to display (info,debug,warn,error) (default "info")
#   -slack-token string
#         [REQUIRED] Slack token used to authenticate with api
#   -version
#         Get current version


# Run Bashbot
bashbot \
  --config-file $BASHBOT_CONFIG_FILEPATH \
  --slack-token $SLACK_TOKEN \
  --log-level info \
  --log-format text
# INFO[2021-07-21T15:19:12-07:00] BashBot Started: 2021-07-21 15:19:12.054219 -0700 PDT m=+0.001068762 
# INFO[2021-07-21T15:19:12-07:00] Bashbot is now connected to slack. Primary trigger: `BashBot`
```


--------------------------------------------------

### Run Bashbot as a docker container

***Requirements***

- docker
- wget

```bash
# Get the sample config.json
wget -O config.json https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.json

# Set environment variables
export SLACK_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
export BASHBOT_CONFIG_FILEPATH=${PWD}/config.json

# Set environment variables and run bashbot
docker run \
  -v ${PWD}/config.json:/bashbot/config.json \
  -v ${PWD}/.env:/bashbot/.env \
  -e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
  -e BASHBOT_ENV_VARS_FILEPATH="/bashbot/.env" \
  -e SLACK_TOKEN=$SLACK_TOKEN \
  -e LOG_LEVEL="info" \
  -e LOG_FORMAT="text" \
  -it mathewfleisch/bashbot:v1.4.6
```

### Run bashbot in kubernetes

***Requirements***

- kubernetes

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: bashbot
  name: bashbot
  namespace: bashbot
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: bashbot
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: bashbot
    spec:
      containers:
      containers:
      -
        env:
          -
            name: BASHBOT_ENV_VARS_FILEPATH
            value: /bashbot/.env
        image: mathewfleisch/bashbot:v1.4.6
        imagePullPolicy: IfNotPresent
        name: bashbot
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        workingDir: /bashbot
        volumeMounts:
        - name: config-json
          mountPath: /bashbot/config.json
        - name: env-vars
          mountPath: /bashbot/.env
      volumes:
        - name: config-json
          hostPath:
            path: /tmp/bashbot-configs/bashbot/config.json
        - name: env-vars
          hostPath:
            path: /tmp/bashbot-configs/bashbot/.env
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30

```


-------------------------------------------------------------------------


***Steps To Prove It's Working***

- Now you should be able to run a few commands in your slack channel ...
- Create a new public channel in your slack called `#bot-test`
- Invite the BashBot into your channel by typing `@BashBot`
- Slackbot should respond with the message: `OK! I’ve invited @BashBot to this channel.`
- Now type `bashbot help`
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.json









-------------------------------------------------------------------------

### config.json
[sample-config.json](sample-config.json)
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
      "command": "cat config.json | jq -r '.tools[] | .trigger' | sort",
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
      "install": "git clone https://$GITHUB_TOKEN@github.com/eaze/bashbot-scripts.git"
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

In this example, a user would type `bashbot list-commands` and that would then run the command `cat config.json | jq -r '.tools[] | .trigger' | sort` which takes no parameters and returns a code block of text from the response. 
```json
{
  "name": "List Commands",
  "description": "List all of the possible commands stored in bashbot",
  "help": "bashbot list-commands",
  "trigger": "list-commands",
  "location": "./",
  "setup": "echo \"\"",
  "command": "cat config.json | jq -r '.tools[] | .trigger' | sort",
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```
#### parameters continued (1 of 2):
There are a few ways to define parameters for a command, but the parameters passed to the bot MUST be white-listed. If the command can be triggered with no parameters, an empty array can be used as in the first example. If the command requires parameters, they can be a hard coded array of strings, or derived from another command. In this example, the hard-coded list of possible parameters is defined (random, question, answer). The `question` action essentially selects a random line in the `--questions-file` text file. The `answer` action does the same as questions, but with a greater-than sign appended to the string. Finally, the `random` action selects both, a random question and random answer from both linked text files.
```json
{
  "name": "Cards Against Humanity",
  "description": "Picks a random question and answer from a list.",
  "help": "bashbot cah [random|question|answer]",
  "trigger": "cah",
  "location": "./vendor/bashbot-scripts",
  "setup": "echo \"\"",
  "command": "./cardsAgainstHumanity.sh --action ${action} --questions-file ../against-humanity/questions.txt --answers-file ../against-humanity/answers.txt",
  "parameters": [{
    "name": "action",
    "allowed": ["random", "question", "answer"]
  }],
  "log": false,
  "ephemeral": false,
  "response": "text",
  "permissions": ["all"]
}
```
#### parameters continued (2 of 2): 
In this example, a list of all 'trigger' values are extracted from the config.json and used as the parameter white-list. When the parameter list can be derived from output of another unix command, it can be "piped" in using the 'source' key. The command must execute without additional stdin input and consist of a newline separated list of values. The command jq is used to parse the json file of all 'trigger' values in a newline separated list.
```json
{
  "name": "Describe Command",
  "description": "Show the json object for a specific command",
  "help": "bashbot describe [command]",
  "trigger": "describe",
  "location": "./scripts",
  "setup": "echo \"\"",
  "command": "./describe-command.sh ../config.json ${command}",
  "parameters": [
    {
      "name": "command",
      "allowed": [],
      "description": "a command to describe ('bashbot list-commands')",
      "source": "cat ../config.json | jq -r '.tools[] | .trigger'"
    }],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```



## Automation (Build/Release)
Included in this repository two github actions are executed on git tags. The [![build-release](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml) action will build multiple go-binaries for each version (linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64) and add them to a github release. The 
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml) action will use the docker plugin, buildx, to build and push a container for amd64/arm64 to docker hub.

```bash
# example semver bump: v1.4.6
git tag v1.4.6
git push origin v1.4.6
```