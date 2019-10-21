# BashBot

BashBot is a white-listed command injection tool for slack. A [config.json](sample-config.json) file defines the possible commands that can be run as well as all of the parameters that can be passed to those commands. This bot uses circleci to build a docker container, that is pushed to AWS ECR and is run in ECS. Sensitive commands can be restricted to specific slack channels. Import other repositories like [bashbot-scripts](https://github.com/eaze/bashbot-scripts) to extend functionality, and reduce clutter in the configuration file.

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

### Sample admin.json
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

### Sample config.json
[sample-config.json](sample-config.json)
```
# Config files are defined as an array of json objects
$ cat sample-config.json | jq '.tools[4]'
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

# name, description and help provide human readable information about the specific command
# trigger:      unique alphanumeric word that represents the command
# location:     absolute or relative path to dependency directory (use "./" for no dependency)
# setup:        command that is run before the main command. (use "echo \"\"" as a default)
# command:      bash command using ${parameter-name} to inject white-listed parameters or environment variables
# parameters:   array of parameter objects. (more detail below)
# log:          define whether the command should be logged in log channel
# ephemeral:    define if the response should be shown to all, or just the user that triggered the command
# response:     [code|text] code displays response in a code block, text displays response as raw text
# permissions:  array of strings. private channel ids to restrict command access to

# parameters continued:
# There are a few ways to define parameters for a command, but the parameters passed to the bot MUST be white-listed. This is an attempt to ensure security of sensitive environment variables. If the command can be triggered with no parameters, an empty array can be used as in the first example. If the command requires parameters, they can be a hard coded array of strings, or derived from another command. In this example, the hard-coded list of possible parameters is defined (random, question, answer). The 'question' action essentially selects a random line in the 'questions-file' text file. The 'answer' action does the same as questions, but with a greater-than sign appended to the string. Finally, the 'random' action selects both, a random question and random answer from both linked text files.
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

# parameters continued:
In this example, a list of slack user ids are derived from a raw api request and used as the parameter white-list.
{
  "name": "Slap User",
  "description": "Slap a specific user with a trout gif",
  "help": "bashbot slap [user]",
  "trigger": "slap",
  "location": "./vendor/bashbot-scripts",
  "setup": "echo \"\"",
  "command": "./giphy.sh slap+trout 10",
  "parameters": [
    {
      "name": "user",
      "allowed": [],
      "description": "tag a slack user",
      "source": "curl -s \"https://slack.com/api/users.list?token=$SLACK_TOKEN\" | jq -r '.members[] | select(.deleted == false) | .id' | sort | sed -e 's/\\(.*\\)/<@\\1>/g'"
    }
  ],
  "log": false,
  "ephemeral": false,
  "response": "text",
  "permissions": ["all"]
}
```

### Sample messages.json
[sample-messages.json](sample-messages.json)

### Setup (bare metal)

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

./start.sh
```
### Setup (Docker)

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

### Setup (ECS)

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

### CircleCi Environment Variables
```
# AWS credentials
AWS_ACCESS_KEY_ID=xxxxxxxxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxx
REMOTE_CONFIG_BUCKET=xxxxxxxxxxxxx
```
