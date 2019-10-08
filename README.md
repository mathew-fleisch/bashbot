# BashBot

BashBot is a whitelisted command injection tool for slack. A [config.json](config.json) file defines the possible commands that can be run as well as all of the parameters that can be passed to those commands. This bot uses circleci to build a docker container, that is pushed to AWS ECR and is run in ECS. Sensitive commands can be restricted to specific slack channels. 

### Example Configuration

  - Fork bashbot to your organization
  - Setup redis t2.micro (note endpoint/port)
  - Create public and private s3 buckets to setup aws and store secrets
  - Setup ecs cluster, task definition, service and repsitory
  - Define a .env file for environment variables save to private bucket
  - Define a config.json and admin.json file and save to private bucket

### CircleCi Environment Variables
```
# AWS credentials
AWS_ACCESS_KEY_ID=xxxxxxxxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxx
REMOTE_CONFIG_BUCKET=xxxxxxxxxxxxx
```

### Required Environment Variables (.env file)
```
# GitHub credentials
GITHUB_USER=xxxxxxxxxxxxx
GITHUB_TOKEN=xxxxxxxxxxxxx

# Public/Private s3 buckets
AWS_PUBLIC_SETUP_URL=xxxxxxxxxxxxx
REMOTE_CONFIG_BUCKET=xxxxxxxxxxxxx

# ECS Information
ECS_REPO=xxxxxxxxxxxxx
ECS_CLUSTER=xxxxxxxxxxxxx
ECS_SERVICE=xxxxxxxxxxxxx
ECS_URL=xxxxxxxxxxxxx
ECS_REGION=xxxxxxxxxxxxx
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
```
{
  "tools": [{
      "name": "BashBot Help",
      "description": "Show this message",
      "help": "bashbot help",
      "trigger": "help",
      "location": "./",
      "setup": "echo \"BashBot is a white-listed command injection tool for slack... written in go. Add this bot to the channel that you wish to carry out commands, and type \\`bashbot help\\` to see this message.\nRun \\`bashbot <command> help\\` to see whitelist of parameters.\nPossible \\`<commands>\\`:\"",
      "command": "echo \"\\`\\`\\`\" && cat config.json | jq -r -c '.tools[] | \"\\(.help) - \\(.description)\"' && echo \"\\`\\`\\`\"",
      "parameters": [],
      "log": false,
      "ephemeral": false,
      "response": "text",
      "permissions": ["all"]
    }
  ],
  "dependencies": [
    {
      "name": "BashBot scripts Scripts",
      "source": "https://$GIT_TOKEN@github.com/eaze/bashbot-scripts.git",
      "install": "git clone ${source}",
      "setup": "echo \"\""
    }
  ]
}

```