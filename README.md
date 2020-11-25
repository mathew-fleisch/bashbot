```
---------------------------------------
 ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|
---------------------------------------
```
BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to trigger bash scripts. A [configuration file](sample-config.json) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls.


--------------------------------------------------

## Installation and setup 

Bashbot can be run as a go binary or as a container and requires an .env file for secrets/environment-variables and a config.json saved in a _git repository_. The .env file will contain a slack token, a git token (for pulling private repositories), and the location of a config.json file. This _git repository_ should exist in your organization/personal-github-account and should be devoted to your configuration of the bot. Bashbot will read from this repository constantly, making it easy to change the configuration without restarting the bot. An s3 bucket can used to store the .env file and referenced via environment variables to pull configuration/secrets for specific bashbot instances.

***Note***

Slack's permissions model has changed and the "[RTM](https://api.slack.com/rtm)" socket connection requires a "classic app" to be configured to get the correct type of token to run Bashbot. After logging into slack, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "legacy bot user" and "Bot User OAuth Access Token" and save the `xoxb-xxxxxxxxx-xxxxxxx` as the environment variable `SLACK_TOKEN` in a `.env` file at bashbot's root.

***Requirements***

- jq
- git
- golang

or

- docker

```bash
# Step 1: Get slack "classic app" bot token
# https://api.slack.com/apps?new_classic_app=1
#
###################################################
#
# Step 2: Create configuration repository (see below for .env and config.json format)
# botname
#  ├── config.json
#  └── .env
# Note: .env will contain secrets and can be stored in s3 by setting these environment variables: 
#   AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, S3_CONFIG_BUCKET
#
###################################################
#
# Step 3a: Run Bashbot locally
curl -sL https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest \
    | jq -r '.tarball_url' \
    | xargs -I {} curl -sL {} -o bashbot.tar.gz \
  && mkdir bashbot \
  && tar -zxvf bashbot.tar.gz -C bashbot --strip-components=1
# or clone via ssh
git clone git@github.com:mathew-fleisch/bashbot.git
# or clone via https
git clone https://github.com/mathew-fleisch/bashbot.git
# Copy .env/config.json to bashbot root and run entrypoint
cp config.json bashbot/. \
  && cp .env bashbot/. \
  && cd bashbot \
  && ./entrypoint.sh
#
###################################################
#
# Step 3b: Run Bashbot via docker (https://hub.docker.com/r/mathewfleisch/bashbot)
docker run \
  -v ${PWD}/config.json:/bashbot/config.json \
  -v ${PWD}/.env:/bashbot/.env \
  -it mathewfleisch/bashbot:v1.2.0
# or by s3 bucket
docker run \
  -e AWS_ACCESS_KEY_ID="xxx"
  -e AWS_SECRET_ACCESS_KEY="xxx"
  -e S3_CONFIG_BUCKET="s3://[PATH-TO-ENV-FILE]"
  -it mathewfleisch/bashbot:v1.2.0
```
```yaml
#
###################################################
#
# Step 3c: Deploy Bashbot via kubernetes
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
      - image: mathewfleisch/bashbot:v1.2.0
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

### .env file

```bash
export SLACK_TOKEN=<xoxb-xxxxxx-xxxxxx>
export GIT_TOKEN=<generate with permissions to READ from repo defined below>
export github_org=<github user or organization>
export github_repo=<repo that will store the bot config>
export github_branch=<config can be pulled from any branch>
export github_filename=<path/to/config.json>
```



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
Included in this repository is a [github action](.github/workflows/release.yaml) that when git tags are pushed, a release is cut, added to github and a docker build is pushed to a target registry, if the proper github secrets are added in the repository settings. 

```bash
# Expected github secrets for automation to execute
REGISTRY_USERNAME
REGISTRY_PASSWORD
REGISTRY_URL=docker.io
REGISTRY_APPNAME=mathewfleisch/bashbot
FORK_OWNER=mathew-fleisch
GIT_TOKEN

# example semver bump: v1.2.4
git tag v1.2.4
git push origin v1.2.4
```