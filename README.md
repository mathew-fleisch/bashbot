# Bashbot

[![Build binaries](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-release.yaml)
[![Build containers](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/build-container.yaml)
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated)

BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to determine when to trigger bash commands. A [configuration file](sample-config.json) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls.

See the [examples](examples) directory for more information about configuring and customizing Bashbot for your team.

See the [Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more information about how to deploy Bashbot in your infrastructure. In this example, a user triggers a Jenkins job using Bashbot and another instance of Bashbot is deployed in a Jenkins job as a gating mechanism. The configuration for the secondary Bashbot could get info about the Jenkins job/host and provides controls to manually decide if the job should pass or fail, at a certain stage in the build. This method of deploying Bashbot gives basic Jenkins controls (trigger, pass, fail) to users in an organization, without giving them access to Jenkins itself. Bashbot commands can be restricted to private channels to limit access within slack.

<img src="https://i.imgur.com/P6IL10y.gif" />

---

## Installation and setup

Bashbot can be run as a go binary or as a container and requires a slack-token, slack-app-token and a config.json. The go binary takes flags to set the slack-token, slack-app-token and path to the config.json file and the container uses environment variables to trigger a go binary by [entrypoint.sh](entrypoint.sh).

***Note about slack-token and slack-app-token*

Slack's permissions model for the "[Socket Mode](https://api.slack.com/apis/connections/sockethttps://api.slack.com/rtm)" socket connection, requires a "classic app" to be configured to get the correct type of token to run Bashbot. After logging into slack via browser, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "legacy bot user", "Bot User OAuth Access Token" and "[App-Level Token](https://api.slack.com/authentication/token-types#app)". Turn on Socket Mode, **subscribe to bot events you want Bashbot to listen to**. Finally, add bashbot to your workspace and invite to a channel. See the [Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more detailed information about how to deploy Bashbot in your infrastructure.

***Quick start***

```bash
# Set `Bot User OAuth Access Token` as SLACK_BOT_TOKEN environment variable
export SLACK_BOT_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
# Set `Bot App Access Token` as SLACK_APP_TOKEN environment variable
export SLACK_APP_TOKEN=xapp-xxxxxxxxx-xxxxxxx

# Get the sample config.json
wget -O config.json https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.json

# Pass environment variable and mount configuration json to run container
docker run --rm \
  --name bashbot \
  -v ${PWD}/config.json:/bashbot/config.json \
  -e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
  -e SLACK_BOT_TOKEN=${SLACK_BOT_TOKEN} \
  -e SLACK_APP_TOKEN=${SLACK_APP_TOKEN} \
  -e LOG_LEVEL="info" \
  -e LOG_FORMAT="text" \
  -it mathewfleisch/bashbot:latest
```

---

To set up a local [KinD cluster](https://kind.sigs.k8s.io/docs/user/quick-start/) run the following commands (prerequisites: [docker](https://docs.docker.com/get-docker/), [helm 3+](https://helm.sh/docs/intro/quickstart/), [kubectl](https://kubernetes.io/docs/tasks/tools/), [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/) ):

```bash
# Copy .env/config.json/.tool-versions files to helm directory (set .env variable values)
cp sample-env-file helm/bashbot/.env
cp sample-config.json helm/bashbot/config.json
cp .tool-versions helm/bashbot/.tool-versions
make kind-test
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
# helm install bashbot helm/bashbot \
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
# ./helm/bashbot/test-deployment.sh
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
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.json

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

