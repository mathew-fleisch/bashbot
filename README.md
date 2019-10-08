# BashBot

BashBot is a whitelisted command injection tool for slack. A [config.json](config.json) file defines the possible commands that can be run as well as all of the parameters that can be passed to those commands. This bot uses circleci to build a docker container, that is pushed to AWS ECR and is run in ECS. Sensitive commands can be restricted to specific slack channels. Import other repositories like [bashbot-scripts](https://github.com/eaze/bashbot-scripts) to extend functionality, and reduce clutter in the configuration file.


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

### Setup (Docker)

  - Fork bashbot to your organization
  - Create public and private s3 buckets to setup aws and store secrets
  - Define a .env file for environment variables save to private bucket and root of bashbot
  - Define a config.json and admin.json file and save to private bucket and root of bashbot

```
# Create .env, config.json and admin.json

# Build Docker Container
docker build -t bashbot .

# Run Docker Container
docker run -it bashbot:latest
```

### Setup (ECS)

  - Fork bashbot to your organization
  - Create public and private s3 buckets to setup aws and store secrets
  - Setup ecs cluster, task definition, service and repository
  - Define a .env file for environment variables save to private bucket
  - Define a config.json and admin.json file and save to private bucket

### CircleCi Environment Variables
```
# AWS credentials
AWS_ACCESS_KEY_ID=xxxxxxxxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxx
REMOTE_CONFIG_BUCKET=xxxxxxxxxxxxx
```
