# Bashbot

[![Build binaries](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml)
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml)
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated)

BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to determine when to trigger bash commands. A [configuration file](sample-config.json) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls. See the [examples](examples) directory, for more information about configuring and customizing Bashbot, for your team.


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
# bashbot-darwin-amd64  v1.5.4

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
#         [REQUIRED] Filepath to config.json file (or environment variable BASHBOT_CONFIG_FILEPATH set)
#   -help
#         Help/usage information
#   -install-vendor-dependencies
#         Cycle through dependencies array in config file to install extra dependencies
#   -log-format string
#         Display logs as json or text (default "text")
#   -log-level string
#         Log level to display (info,debug,warn,error) (default "info")
#   -send-message-channel string
#         Send stand-alone slack message to this channel (requires -send-message-text)
#   -send-message-ephemeral
#         Send stand-alone ephemeral slack message to a specific user (requires -send-message-channel -send-message-text and -send-message-user)
#   -send-message-text string
#         Send stand-alone slack message (requires -send-message-channel)
#   -send-message-user string
#         Send stand-alone ephemeral slack message to this slack user (requires -send-message-channel -send-message-text and -send-message-ephemeral)
#   -slack-token string
#         [REQUIRED] Slack token used to authenticate with api (or environment variable SLACK_TOKEN set)
#   -version
#         Get current version

# Run Bashbot
bashbot \
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
  -it mathewfleisch/bashbot:v1.5.4
```

#### Build bashbot docker container

```bash
# Build bashbot container locally
docker build -t bashbot-local .

# Run local bashbot container
docker run \
  -v ${BASHBOT_CONFIG_FILEPATH}:/bashbot/config.json \
  -e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
  -e SLACK_TOKEN=${SLACK_TOKEN} \
  -e LOG_LEVEL="info" \
  -e LOG_FORMAT="text" \
  --name bashbot --rm \
  -it bashbot-local:latest

# Exec into bashbot container
docker exec -it $(docker ps -aqf "name=bashbot") bash

# Remove existing bashbot container
docker rm $(docker ps -aqf "name=bashbot")
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
    type: Recreate
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
        image: mathewfleisch/bashbot:v1.5.4
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
      terminationGracePeriodSeconds: 0

```


-------------------------------------------------------------------------


***Steps To Prove It's Working***

- Now you should be able to run a few commands in your slack channel ...
- Create a new public channel in your slack called `#bot-test`
- Invite the BashBot into your channel by typing `@BashBot`
- Slackbot should respond with the message: `OK! I’ve invited @BashBot to this channel.`
- Now type `bashbot help`
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.json



## Configuration (config.json)

The configuration of bashbot is loaded into the go-binary as a json file. More about the config.json syntax can be found in the [examples](examples) directory.


## Automation (Build/Release)
Included in this repository two github actions are executed on git tags. The [![build-release](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml) action will build multiple go-binaries for each version (linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64) and add them to a github release. The 
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml) action will use the docker plugin, buildx, to build and push a container for amd64/arm64 to docker hub.

```bash
# example semver bump: v1.5.4
git tag v1.5.4
git push origin v1.5.4
```