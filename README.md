# Bashbot

[![Build binaries](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml)
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml)
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated)

BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to determine when to trigger bash commands. A [configuration file](sample-config.json) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls. 

See the [examples](examples) directory for more information about configuring and customizing Bashbot for your team.

See the [Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more information about how to deploy Bashbot in your infrastructure. In this example, a user triggers a Jenkins job using Bashbot and another instance of Bashbot is deployed in a Jenkins job as a gating mechanism. The configuration for the secondary Bashbot could get info about the Jenkins job/host and provides controls to manually decide if the job should pass or fail, at a certain stage in the build. This method of deploying Bashbot gives basic Jenkins controls (trigger, pass, fail) to users in an organization, without giving them access to Jenkins itself. Bashbot commands can be restricted to private channels to limit access within slack.

<img src="https://i.imgur.com/P6IL10y.gif" />

--------------------------------------------------

## Installation and setup

Bashbot can be run as a go binary or as a container and requires a slack-token and a config.json. The go binary takes flags to set the slack-token and path to the config.json file and the container uses environment variables to trigger a go binary by [entrypoint.sh](entrypoint.sh).

***Note about slack-token***

Slack's permissions model for the "[Real-Time-Messaging (RTM)](https://api.slack.com/rtm)" socket connection, requires a "classic app" to be configured to get the correct type of token to run Bashbot. After logging into slack via browser, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "legacy bot user" and "Bot User OAuth Access Token." Finally, add bashbot to your workspace and invite to a channel. See the [Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more detailed information about how to deploy Bashbot in your infrastructure.


***Quick start***

```bash
# Set `Bot User OAuth Access Token` as SLACK_TOKEN environment variable
export SLACK_TOKEN=xoxb-xxxxxxxxx-xxxxxxx

# Get the sample config.json
wget -O config.json https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.json

# Pass environment variable and mount configuration json to run container
docker run --rm \
   --name bashbot \
   -v ${PWD}/config.json:/bashbot/config.json \
   -e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
   -e SLACK_TOKEN=${SLACK_TOKEN} \
   -e LOG_LEVEL="info" \
   -e LOG_FORMAT="text" \
   -it mathewfleisch/bashbot:latest
```

-------------------------------------------------------------------------

 - [Bashbot Setup/Deployment Examples](https://github.com/mathew-fleisch/bashbot-example)
    - Setup (one time)
      - [Setup Step 0: Make the slack app and get a token](https://github.com/mathew-fleisch/bashbot-example#setup-step-0-make-the-slack-app-and-get-a-token)
      - [Setup Step 1: Fork this repository](https://github.com/mathew-fleisch/bashbot-example#setup-step-1-fork-this-repository)
      - [Setup Step 2: Create deploy key](https://github.com/mathew-fleisch/bashbot-example#setup-step-2-create-deploy-key) (handy for private "forks")
      - [Setup Step 3: Upload public deploy key](https://github.com/mathew-fleisch/bashbot-example#setup-step-3-upload-public-deploy-key)
      - [Setup Step 4: Save deploy key as github secret](https://github.com/mathew-fleisch/bashbot-example#setup-step-4-save-deploy-key-as-github-secret) (optional: used in github-actions)

    - Running Bashbot Locally
      - [***Go-Binary***](https://github.com/mathew-fleisch/bashbot-example#run-bashbot-locally-as-go-binary)
      - [***Docker***](https://github.com/mathew-fleisch/bashbot-example#run-bashbot-locally-from-docker)

    - Deploy Bashbot
      - [Kubernetes](https://github.com/mathew-fleisch/bashbot-example#run-bashbot-in-kubernetes)

-------------------------------------------------------------------------

 - [Configuration Examples](examples)
    1. Simple call/response
        - [Ping/Pong](examples/ping)
    2. Run bash script
        - [Get Caller Information](examples/info)
        - [Get asdf Plugins](examples/asdf)
        - [Get add-example Plugins](examples/add-example)
        - [Get remove-example Plugins](examples/remove-example)
        - [Get list-examples Plugins](examples/list-examples)
    3. Run golang script
        - [Get Running Version](examples/version)
    4. Parse json via jq
        - [Show Help Dialog](examples/help)
        - [List Available Commands](examples/list)
        - [Describe Command](examples/describe)
    5. Curl/wget
        - [Get Latest Bashbot Release](examples/latest-release)
        - [Get File From Repository](examples/get-file-from-repo)
        - [Trigger Github Action](examples/trigger-github-action)
        - [Get Air Quality Index By Zip](examples/aqi)
    6. [Send message from github action](examples/#send-message-from-github-action)



***Steps To Prove It's Working***

- Now you should be able to run a few commands in your slack channel ...
- Create a new public channel in your slack called `#bot-test`
- Invite the BashBot into your channel by typing `@BashBot`
- Slackbot should respond with the message: `OK! I’ve invited @BashBot to this channel.`
- Now type `bashbot help`
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.json



-------------------------------------------------------------------------

## Automation (Build/Release)
Included in this repository two github actions are executed on git tags. The [![build-release](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml) action will build multiple go-binaries for each version (linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64) and add them to a github release. The
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml) action will use the docker plugin, buildx, to build and push a container for amd64/arm64 to docker hub.

```bash
# example semver bump: v1.6.15
git tag v1.6.15
git push origin v1.6.15
```

There are also automated anchore container scans and codeql static analysis done on every push to the main branch.

