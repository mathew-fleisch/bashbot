# Bashbot

[![Release](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml) |
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated) | [ghcr](https://github.com/mathew-fleisch/bashbot/pkgs/container/bashbot)

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

### Quick Start: Helm Install

To install bashbot from the public helm repo, rather than with the source, the slack app/bot tokens are required to be saved as kubernetes secret, and local versions of the config.yaml and .tool-versions file.

```bash
# Get sample-config.yaml and .tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/.tool-versions -q -O ${PWD}/.tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml -q -O ${PWD}/config.yaml

# Add the public helm repo
helm repo add bashbot https://mathew-fleisch.github.io/bashbot

# Define your bashbot instance and namespace names
export BOTNAME=bashbot
export NAMESPACE=bashbot
helm install \
  --debug --wait \
  --namespace ${NAMESPACE} \
  --set namespace=${NAMESPACE} \
  --set botname=${BOTNAME} \
  --set-file '\.tool-versions'=${PWD}/.tool-versions \
  --set-file 'config\.yaml'=${PWD}/config.yaml \
  --repo https://mathew-fleisch.github.io/bashbot \
  bashbot bashbot

```

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
helm 3.10.3
kubectl 1.26.0
kubectx 0.9.4
kustomize 4.5.7
```

To set up a local [KinD cluster](https://kind.sigs.k8s.io/docs/user/quick-start/) run the following commands:

```bash
# Copy .env, config.yaml and .tool-versions sample files to helm directory and replace with custom values.
cp sample-env-file .env
cp sample-config.yaml config.yaml
make test-kind

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

The go-binaries are saved as [release artifacts](https://github.com/mathew-fleisch/bashbot/releases) and can be downloaded with curl/wget or by the [asdf version manager](https://asdf-vm.com/) `asdf plugin add bashbot && asdf install bashbot latest && asdf global bashbot latest` (skip to `bashbot version` if installed with asdf in the following instructions)

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
export latest_bashbot_version=$(curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)

# Remove any existing bashbot binaries
rm -rf /usr/local/bin/bashbot || true

# Get correct binary for host machine and place in user's path
wget -qO /usr/local/bin/bashbot https://github.com/mathew-fleisch/bashbot/releases/download/${latest_bashbot_version}/bashbot-${os}-${arch}

# Make bashbot binary executable
chmod +x /usr/local/bin/bashbot

# Disable security check on macs (only run once)
test "${os}" == "darwin" && xattr -d com.apple.quarantine /usr/local/bin/bashbot

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
bashbot install-dependencies \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"

# Run Bashbot binary passing the config file and the Slack token
bashbot run \
  --log-level "$LOG_LEVEL" \
  --log-format "$LOG_FORMAT"
```

---

The [Examples](examples) directory of this repository, has many commands used in automated tests, and can illustrate how bashbot can be used to trigger automation, or as automation itself, by leveraging secrets from the host running bashbot. For instance, one command might use an api token to curl a third-party api, and return the response back to the slack user, once triggered ([aqi example](examples/aqi)). If deployed in a kubernetes cluster, with a service-account to access the kube-api, bashbot can allow slack users to execute (hopefully curated) kubectl/helm commands, for an SRE/Ops focused deployment (get/describe/delete pods/deployments/secrets etc).

1. Simple call/response
    - [Ping/Pong](examples/ping)
2. Run bash script
    - [Get Caller Information](examples/info)
    - [Get asdf Plugins](examples/asdf)
    - [Get list-examples Plugins](examples/list-examples)
3. Run golang script
    - [Get Running Version](examples/version)
4. Parse json via jq/yq
    - [Show Help Dialog](examples/help)
    - [List Available Commands](examples/list)
    - [Describe Command](examples/describe)
5. Curl/wget
    - [Get Latest Bashbot Release](examples/latest-release)
    - [Get File From Repository](examples/get-file-from-repo)
    - [Get Air Quality Index By Zip](examples/aqi)
6. In Github Actions (or other CI)
    - [Trigger Github Action](examples/trigger-github-action)
    - [Trigger Gated Github Action](examples/trigger-github-action)
7. Kubernetes
    - [kubectl cluster-info](kubernetes#kubectl-cluster-info)
    - [kubectl get pod](examples/kubernetes#kubectl-get-pod)
    - [kubectl -n [namespace] delete pod [pod-name]](examples/kubernetes#kubectl--n-namespace-delete-pod-pod-name)
    - [kubectl -n [namespace] decribe pod [podname]](examples/kubernetes#kubectl--n-namespace-decribe-pod-podname)
    - [kubectl -n [namespace] logs -f [pod-name]](examples/kubernetes#kubectl--n-namespace-logs--f-pod-name)

***Steps To Prove It's Working***

- Now you should be able to run a few commands in your slack channel ...
- Create a new public channel in your slack called `#bot-test`
- Invite the BashBot into your channel by typing `@BashBot`
- Slackbot should respond with the message: `OK! I’ve invited @BashBot to this channel.`
- Now type `!bashbot help`
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.yaml

---

## Automation (Build/Release)

Included in this repository two github actions are executed on git tags. The [![release](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml) action will build multiple go-binaries for each version (linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64) and add them to a github release and will use the docker plugin, buildx, to build and push a container for amd64/arm64 to docker hub and github container registry.

```bash
# example semver bump: v2.0.0
git tag v2.0.0
git push origin v2.0.0
```

On [pull requests to the main branch](.github/workflows/pr.yaml), four jobs are run on every commit:

- linting and unit tests are run under the `unit_tests` job
- a container is built and scanned by the [anchore](https://anchore.com/) container scanning tool
- the golang code is analyzed by [codeql](https://codeql.github.com/) SAST tool
- a container is built and deployed in a kind cluster, to run automated tests, to maintain/verify basic functionality (see [Makefile](Makefile) target `make help` for more information)

### Makefile

run `make help` for a full list of targets. Any of these environment variables can be overridden by exporting a new value before running any makefile target.

```makefile
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION?=$(shell make version)
LATEST_VERSION?=$(shell curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)
BINARY?=bin/bashbot
SRC_LOCATION?=main.go
BASHBOT_LOG_LEVEL?=info
BASHBOT_LOG_TYPE?=text
TESTING_CHANNEL?=C034FNXS3FA
ADMIN_CHANNEL?=GPFMM5MD2
NAMESPACE?=bashbot
BOTNAME?=bashbot
HELM_CONFIG_YAML?=$(PWD)/config.yaml
HELM_TOOL_VERSIONS?=$(PWD)/.tool-versions
HELM_ENV?=${PWD}/.env
```

Note: [yq](https://github.com/mikefarah/yq) is a dependency of running many makefile targets and can be installed with the binary: `wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -q -O /usr/local/bin/yq && chmod +x /usr/local/bin/yq`
