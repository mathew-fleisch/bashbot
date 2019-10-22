# BashBot

BashBot is a white-listed command injection tool for slack. A [config.json](sample-config.json) file defines the possible commands that can be run as well as all of the parameters that can be passed to those commands. This bot uses circleci to build a docker container, that is pushed to AWS ECR and is run in ECS. Sensitive commands can be restricted to specific slack channels. Import other repositories like [bashbot-scripts](https://github.com/eaze/bashbot-scripts) to extend functionality, and reduce clutter in the configuration file.

## Table of Contents

- [Installation](#Installation%20and%20setup)
  * [Slack App Setup](#Slack%20App%20Setup)
  * [Bare Metal](#Bare%20Metal%20Setup)
  * [Docker](#Docker%20Setup)
  * [ECS](#ECS%20Setup)
- [Sample Env File](#Sample%20.env%20file)
- [Sample admin.json](#Sample%20admin.json)
- [Sample config.json](#Sample%20config.json)
- [Sample messages.json](#Sample%20messages.json)
- [CircleCi Environment Variables](#CircleCi%20Environment%20Variables)

## Installation and setup 
We have listed 3 different ways to install and get this up and running! Sample .env, admin.json, config.json and sample messages below. When setting up your s3 and ecs cluster make sure they are in the same region.

### Slack App Setup

- Copy "Bot User OAuth Access Token" to .env file

### Bare Metal Setup

```
# Log in as root
sudo -i

# Install dependencies
apt update
apt upgrade -y
apt install -y zip wget iputils-ping curl jq build-essential libssl-dev ssh python python-pip python3 python3-pip openssl file libgcrypt-dev git redis-server sudo build-essential libssl-dev awscli vim

# Clone bashbot
git clone https://github.com/eaze/bashbot.git /bashbot

# Create .env file
touch /bashbot/.env
# add secrets/tokens

# Copy Sample Config
cp /bashbot/sample-config.json /bashbot/config.json

# Copy Sample Messages Config
cp /bashbot/sample-messages.json /bashbot/messages.json

# Create Admin Config
touch /bashbot/admin.json
# add personal user id and channel id for public/private channels

# Install Go runtime (version 1.12 at least)
wget https://dl.google.com/go/go1.12.12.linux-amd64.tar.gz
tar xvf go1.12.12.linux-amd64.tar.gz
sudo mv go /usr/local
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.bashrc
source ~/.bashrc

./start.sh
```

----------------------------------------------------------------

### Docker Setup

  - clone bashbot locally
  - Create public and private s3 buckets to setup aws and store secrets
  - Define a .env file for environment variables save to private bucket and root of bashbot
  - Define a config.json, messages.json and admin.json file and save to private bucket and root of bashbot

```
# Create/modify .env, config.json, messages.json, admin.json and `push-configs` to config s3 bucket
./bashbot.sh --action push-configs --config-bucket [bucket-path] 

# Run build command
./bashbot.sh --action build-docker --config-bucket [bucket-path]

# Run Docker Container
docker run -it bashbot:latest
```

----------------------------------------------------------------


### ECS Setup

  - clone bashbot locally
  - Create public and private s3 buckets to setup aws and store secrets
  - Setup ecs cluster, task definition, service and repository
  - Define a .env file for environment variables save to private bucket
  - Define a config.json, messages.json and admin.json file and save to private bucket

```
# Create/modify .env, config.json, messages.json, admin.json and `push-configs` to config s3 bucket
./bashbot.sh --action push-configs --config-bucket [bucket-path] 

# Run build command through circleci job
./bashbot.sh --action build-ecs --config-bucket [bucket-path] --circle-token [circleci-token] --circle-project [circleci-project]
```

---------------------------------------------------------------- 

### Sample .env file
```
# GitHub credentials
export GITHUB_USER=xxxxxxxxxxxxx
export GITHUB_TOKEN=xxxxxxxxxxxxx

# Slack Token
export SLACK_TOKEN=xxxxxxxxxxxxx

# Public/Private s3 buckets
export AWS_PUBLIC_SETUP_URL=xxxxxxxxxxxxx
export REMOTE_CONFIG_BUCKET=xxxxxxxxxxxxx

# ECS Information
export ECS_REPO=xxxxxxxxxxxxx
export ECS_CLUSTER=xxxxxxxxxxxxx
export ECS_SERVICE=xxxxxxxxxxxxx
export ECS_URL=xxxxxxxxxxxxx
export ECS_REGION=xxxxxxxxxxxxx
```

----------------------------------------------------------------

### admin.json
```
{
  "admins": [{
    "trigger": "bashbot",
    "appName": "BashBot",
    "userIds": ["admin-user-id"],
    "privateChannelId": "private-slack-channel-id",
    "logChannelId": "public-log-slack-channel-id"
  }]
}
```

----------------------------------------------------------------

### messages.json
[sample-messages.json](sample-messages.json)

----------------------------------------------------------------


### config.json
[sample-config.json](sample-config.json)
The config.json file is defined as an array of json objects keyed by 'tools' and 'dependencies.' The dependencies section defines any resources that need to be downloaded or cloned from a repository before execution of commands. The following is a simplified example of a config.json file:

```
{
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
      "source": "https://$GITHUB_TOKEN@github.com/eaze/bashbot-scripts.git",
      "install": "git clone ${source}",
      "setup": "echo \"\""
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
```
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
```
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
```
{
  "name": "Describe Command",
  "description": "Show the json object for a specific command",
  "help": "bashbot describe [command]",
  "trigger": "describe",
  "location": "./vendor/bashbot-scripts",
  "setup": "echo \"\"",
  "command": "./describeCommand.sh ../../config.json ${command}",
  "parameters": [
    {
      "name": "command",
      "allowed": [],
      "description": "a command to describe ('bashbot list-commands')",
      "source": "cat ../../config.json | jq -r '.tools[] | .trigger'"
    }],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```

----------------------------------------------------------------


### CircleCi Environment Variables
```
# AWS credentials
AWS_ACCESS_KEY_ID=xxxxxxxxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxx
REMOTE_CONFIG_BUCKET=xxxxxxxxxxxxx
```
