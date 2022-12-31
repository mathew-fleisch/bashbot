# Bashbot

[![Build binaries](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml)
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml)
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated)

BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to determine when to trigger bash commands. A [configuration file](sample-config.yaml) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls.

See the [examples](examples) directory for more information about configuring and customizing Bashbot for your team.

See the [Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more information about how to deploy Bashbot in your infrastructure. In this example, a user triggers a Jenkins job using Bashbot and another instance of Bashbot is deployed in a Jenkins job as a gating mechanism. The configuration for the secondary Bashbot could get info about the Jenkins job/host and provides controls to manually decide if the job should pass or fail, at a certain stage in the build. This method of deploying Bashbot gives basic Jenkins controls (trigger, pass, fail) to users in an organization, without giving them access to Jenkins itself. Bashbot commands can be restricted to private channels to limit access within slack.

<img src="https://i.imgur.com/P6IL10y.gif" />

---

## Installation and setup

Bashbot can be run as a go binary or as a container and requires a slack-bot-token, slack-app-token and a config.yaml. The go binary takes flags to set the slack-bot-token, slack-app-token and path to the config.yaml file and the container uses environment variables to trigger a go binary by [entrypoint.sh](entrypoint.sh).

***Note about slack-bot-token and slack-app-token*

Bashbot uses the Slack API's "[Socket Mode](https://api.slack.com/apis/connections/socket)" to connect to the slack servers over a socket connection and uses ***no webhooks/ingress*** to trigger commands. Bashbot "subscribes" to message and emoji events and determins if a command should be executed and what command should be executed by parsing the data in each event. To run Bashbot, you must have a "[Bot User OAuth Token](https://api.slack.com/authentication/token-types#bot)" and an "[App-Level Token](https://api.slack.com/authentication/token-types#app)".

- Click "Create New App" from the [Slack Apps](https://api.slack.com/apps) page and follow the "From Scratch" prompts to give your instance of bashbot a unique name and a workspace for it to be installed in
- The "Basic Information" page gives controls to set a profile picture toward the bottom (make sure to save any changes)
- Enable "Socket Mode" from the "Socket Mode" page and add the default scopes `conversations.read` and note the "[App-Level Token](https://api.slack.com/authentication/token-types#app)" that is generated to save in the .env file as `SLACK_APP_TOKEN`
- Enable events from the "Event Subscriptions" page and add the following bot event subscriptions and save changes
  - `app_mention`
  - `message.channels`
  - `message.groups`
  - `reaction_added`
  - `reaction_removed`
- From the "OAuth & Permissions" page, after setting the following "scopes" or permissions, install Bashbot in your workspace (this will require administrator approval of your slack workspace) and note the "[Bot User OAuth Token](https://api.slack.com/authentication/token-types#bot)" to save in the .env file as `SLACK_BOT_TOKEN`
  - `app_mentions:read`
  - `channels:history`
  - `channels:read`
  - `chat:write`
  - `files:write`
  - `groups:history`
  - `groups:read`
  - `incoming-webhook`
  - `reactions:read`
  - `reactions:write`
  - `users:read`

Note: Prior to version 2, Bashbot required a "classic app" to be configured to get the correct type of token to connect to Slack. After logging into Slack via browser, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "legacy bot user", "Bot User OAuth Access Token," add bashbot to your workspace, and invite to a channel. See the [Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more detailed information about how to deploy Bashbot in your infrastructure.

### Quick Start: KinD cluster

***KinD Prerequisites***

- [docker](https://docs.docker.com/get-docker/)
- [helm 3+](https://helm.sh/docs/intro/quickstart/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/)

Note: you can use the [asdf version manager](https://asdf-vm.com/) to install dependencies used in the makefile with the following .tool-versions

```text
action-validator 0.1.2
dockle 0.4.3
golangci-lint 1.44.2
helm 3.8.1
kind 0.12.0
yq 4.22.1
```

To set up a local [KinD cluster](https://kind.sigs.k8s.io/docs/user/quick-start/) run the following commands:

```bash
# Copy .env, config.yaml and .tool-versions sample files to helm directory and replace with custom values.
cp sample-env-file charts/bashbot/.env
cp sample-config.yaml charts/bashbot/config.yaml
cp .tool-versions charts/bashbot/.tool-versions
make test-kind
# docker build -t bashbot:local .
# [+] Building 29.8s (17/17) FINISHED  
# ...
#  => => writing image sha256:9d5d360ad6de055ae2f5f9ef859c86b10b8c3613ce078d354f255e3efa8ec000  0.0s
#  => => naming to docker.io/library/bashbot:local  0.0s
# kind create cluster
# Creating cluster "kind" ...
#  ✓ Ensuring node image (kindest/node:v1.21.1) 🖼 
#  ✓ Preparing nodes 📦  
#  ✓ Writing configuration 📜 
#  ✓ Starting control-plane 🕹️ 
#  ✓ Installing CNI 🔌 
#  ✓ Installing StorageClass 💾 
# Set kubectl context to "kind-kind"
# You can now use your cluster with:
# kubectl cluster-info --context kind-kind
# Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
# kind load docker-image bashbot:local
# Image: "bashbot:local" with ID "sha256:9d5d360ad6de055ae2f5f9ef859c86b10b8c3613ce078d354f255e3efa8ec000" not yet present on node "kind-control-plane", loading...
# helm install bashbot charts/bashbot \
#  --namespace bashbot \
#  --create-namespace \
#  --set image.repository=bashbot \
#  --set image.tag=local
# NAME: bashbot
# LAST DEPLOYED: Fri Feb 25 13:07:13 2022
# NAMESPACE: bashbot
# STATUS: deployed
# REVISION: 1
# TEST SUITE: None
# Waiting for bashbot to come up...
# sleep 1
# ./charts/bashbot/test-deployment.sh
# Deployment not found or not ready. 20 more attempts...
# Deployment not found or not ready. 19 more attempts...
# Deployment not found or not ready. 18 more attempts...
# Deployment not found or not ready. 17 more attempts...
# Deployment not found or not ready. 16 more attempts...
# Bashbot deployment confirmed!
# NAME  READY  UP-TO-DATE  AVAILABLE  AGE
# bashbot  1/1  1  1  17s
# Helm test complete!

# At this point you should have a KinD cluster running with a bashbot deployment inside
kubectl cluster-info
helm --namespace bashbot list
# Use these basic kubectl commands (one at a time) to verify (use `-o yaml` for more detailed information)
kubectl --namespace bashbot get deployments
kubectl --namespace bashbot get pods
kubectl --namespace bashbot get configmaps
kubectl --namespace bashbot get serviceaccounts
kubectl --namespace bashbot get clusterrolebinding
kubectl --namespace bashbot get secrets
```

---

### Quick Start: Docker

```bash
# Set `Bot User OAuth Access Token` as SLACK_BOT_TOKEN environment variable
export SLACK_BOT_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
# Set `Bot App Access Token` as SLACK_APP_TOKEN environment variable
export SLACK_APP_TOKEN=xapp-xxxxxxxxx-xxxxxxx

# Get the sample config.yaml
wget -O config.yaml https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml

# Pass environment variable and mount configuration yaml to run container
docker run --rm \
  --name bashbot \
  -v ${PWD}/config.yaml:/bashbot/config.yaml \
  -e BASHBOT_CONFIG_FILEPATH="/bashbot/config.yaml" \
  -e SLACK_BOT_TOKEN=${SLACK_BOT_TOKEN} \
  -e SLACK_APP_TOKEN=${SLACK_APP_TOKEN} \
  -e LOG_LEVEL="info" \
  -e LOG_FORMAT="text" \
  -it mathewfleisch/bashbot:latest
```

---

### Quick Start: Go Binary

```bash
# Set a path to your local configuration file
export BASH_CONFIG_FILEPATH=${PWD}/bashbot/config.yaml
# Set `Bot User OAuth Access Token` as SLACK_BOT_TOKEN environment variable
export SLACK_BOT_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
# Set `Bot App Access Token` as SLACK_APP_TOKEN environment variable
export SLACK_APP_TOKEN=xapp-xxxxxxxxx-xxxxxxx

# ----------- Install binary -------------- #

# os: linux, darwin
export os=$(uname | tr '[:upper:]' '[:lower:]')

# arch: amd64, arm64
export arch=amd64
test "$(uname -m)" == "aarch64" && export arch=arm64

# Latest bashbot version/tag
export latest=$(curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)

# Remove any existing bashbot binaries
rm -rf /usr/local/bin/bashbot || true

# Get correct binary for host machine and place in user's path
wget -qO /usr/local/bin/bashbot https://github.com/mathew-fleisch/bashbot/releases/download/${latest}/bashbot-${os}-${arch}

# Make bashbot binary executable
chmod +x /usr/local/bin/bashbot

# To verify installation run version or help commands
bashbot version
# bashbot-darwin-amd64    v2.0.0

bashbot --help
#  ____            _     ____        _   
# |  _ \          | |   |  _ \      | |  
# | |_) | __ _ ___| |__ | |_) | ___ | |_ 
# |  _ < / _' / __| '_ \|  _ < / _ \| __|
# | |_) | (_| \__ \ | | | |_) | (_) | |_ 
# |____/ \__,_|___/_| |_|____/ \___/ \__|
# Bashbot is a slack bot, written in golang, that can be configured
# to run bash commands or scripts based on a configuration file.

# Usage:
#   bashbot [command]

# Available Commands:
#   completion           Generate the autocompletion script for the specified shell
#   help                 Help about any command
#   install-dependencies Cycle through dependencies array in config file to install extra dependencies
#   run                  Run bashbot
#   send-message         Send stand-alone slack message
#   version              Get current version

# Flags:
#       --config-file string       [REQUIRED] Filepath to config.yaml file (or environment variable BASHBOT_CONFIG_FILEPATH set)
#   -h, --help                     help for bashbot
#       --log-format string        Display logs as json or text (default "text")
#       --log-level string         Log elevel to display (info,debug,warn,error) (default "info")
#       --slack-app-token string   [REQUIRED] Slack app token used to authenticate with api (or environment variable SLACK_APP_TOKEN set)
#       --slack-bot-token string   [REQUIRED] Slack bot token used to authenticate with api (or environment variable SLACK_BOT_TOKEN set)
#   -t, --toggle                   Help message for toggle

# Use "bashbot [command] --help" for more information about a command.


# If the log-level doesn't exist, set it to 'info'
LOG_LEVEL=${LOG_LEVEL:-info}
# If the log-format doesn't exist, set it to 'text'
LOG_FORMAT=${LOG_FORMAT:-text}

# Run install-vendor-dependencies path
bashbot --install-vendor-dependencies \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"

# Run Bashbot binary passing the config file and the Slack token
bashbot \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"
```


---

- [Bashbot Setup/Deployment Examples](https://github.com/mathew-fleisch/bashbot-example)
  - Setup (one time)

  - [Setup Step 0: Make the slack app and get the bot & app token](https://github.com/mathew-fleisch/bashbot-example#setup-step-0-make-the-slack-app-and-get-a-token)
  - [Setup Step 1: Fork this repository](https://github.com/mathew-fleisch/bashbot-example#setup-step-1-fork-this-repository)
  - [Setup Step 2: Create deploy key](https://github.com/mathew-fleisch/bashbot-example#setup-step-2-create-deploy-key) (handy for private "forks")
  - [Setup Step 3: Upload public deploy key](https://github.com/mathew-fleisch/bashbot-example#setup-step-3-upload-public-deploy-key)
  - [Setup Step 4: Save deploy key as github secret](https://github.com/mathew-fleisch/bashbot-example#setup-step-4-save-deploy-key-as-github-secret) (optional: used in github-actions)
  - Running Bashbot Locally

  - [***Go-Binary***](https://github.com/mathew-fleisch/bashbot-example#run-bashbot-locally-as-go-binary)
  - [***Docker***](https://github.com/mathew-fleisch/bashbot-example#run-bashbot-locally-from-docker)
  - Deploy Bashbot

  - [Kubernetes](https://github.com/mathew-fleisch/bashbot-example#run-bashbot-in-kubernetes)

---

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
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.yaml

---

## Automation (Build/Release)

Included in this repository two github actions are executed on git tags. The [![build-release](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml) action will build multiple go-binaries for each version (linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64) and add them to a github release. The
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml) action will use the docker plugin, buildx, to build and push a container for amd64/arm64 to docker hub.

```bash
# example semver bump: v1.8.0
git tag v1.8.0
git push origin v1.8.0
```

There are also automated anchore container scans and codeql static analysis done on every push to the main branch.
