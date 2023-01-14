# Examples and Documentation

- [Installation and Deployments](#installation-and-deployments)
  - [binary](#binary) - Download and install go-binary from github release
  - [binary-asdf](#binary-asdf) - Download and install go-binary from asdf version manager
  - [docker](#docker) - Run a bashbot docker container from dockerhub/ghcr
  - [kubernetes](#kubernetes) - Deploy a raw kubernetes deployment manifest
  - [helm](#helm) - Install bashbot via helm chart
  - [kind](#kind) - Install bashbot via helm chart in a local KinD cluster
  - [argocd](#argocd) - Install bashbot via argocd application pointing to helm chart
- [config.yaml](#configyaml) - Example bashbot commands
  1. Simple call/response
      - [Ping/Pong](ping)
  2. Run bash script
      - [Get Caller Information](info)
      - [Get asdf Plugins](asdf)
      - [Get list-examples Plugins](list-examples)
  3. Run golang script
      - [Get Running Version](version)
  4. Parse json via jq/yq
      - [Show Help Dialog](help)
      - [List Available Commands](list)
      - [Describe Command](describe)
  5. Curl/wget
      - [Get Latest Bashbot Release](latest-release)
      - [Get File From Repository](get-file-from-repo)
      - [Get Air Quality Index By Zip](aqi)
  6. In Github Actions (or other CI)
      - [Trigger Github Action](trigger-github-action)
      - [Trigger Gated Github Action](trigger-github-action)
  7. Kubernetes
      - [kubectl cluster-info](kubernetes#kubectl-cluster-info)
      - [kubectl get pod](kubernetes#kubectl-get-pod)
      - [kubectl -n [namespace] delete pod [pod-name]](kubernetes#kubectl--n-namespace-delete-pod-pod-name)
      - [kubectl -n [namespace] decribe pod [podname]](kubernetes#kubectl--n-namespace-decribe-pod-podname)
      - [kubectl -n [namespace] logs -f [pod-name]](kubernetes#kubectl--n-namespace-logs--f-pod-name)

---

## Installation and Deployments

Bashbot can be deployed in a number of different ways, to leverage the host's environment-variables and network. The [source of bashbot](https://github.com/mathew-fleisch/bashbot) has a makefile to aid in automated testing, that can also be used to run bashbot locally. Clone the source, and run `make help` for an updated list of options, or use this read-me to run commands manually.

**Note:** Prior to version 2, Bashbot required a "classic app" to be configured to get the correct type of token to connect to Slack. After logging into Slack via browser, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "legacy bot user", "Bot User OAuth Access Token," add bashbot to your workspace, and invite to a channel. See the [(archived) Setup/Deployment Examples Repository](https://github.com/mathew-fleisch/bashbot-example) for more detailed information about how to deploy Bashbot (v1) in your infrastructure.

---

### Binary

Bashbot is written in golang, and has release automation to cross compile go-binaries for linux/darwin and amd64/arm64. They can be downloaded from the [releases page](https://github.com/mathew-fleisch/bashbot/releases) and run from the command line with the slack-app-token and slack-bot-token saved as environment variables.

```bash
# Download the default configuration files
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/.tool-versions -q -O ${PWD}/.tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml -q -O ${PWD}/config.yaml
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-env-file -q -O ${PWD}/.env

# Set a path to your local configuration file
export BASH_CONFIG_FILEPATH=${PWD}/config.yaml
# Set `Bot User OAuth Access Token` as SLACK_BOT_TOKEN environment variable
export SLACK_BOT_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
# Set `Bot App Access Token` as SLACK_APP_TOKEN environment variable
export SLACK_APP_TOKEN=xapp-xxxxxxxxx-xxxxxxx
# Set log-level (default: info)
export LOG_LEVEL=info
# Set log-format (default: text) [text|json]
export LOG_FORMAT=text
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
# bashbot-darwin-amd64    v2.0.2

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

### Binary asdf

Bashbot can be installed via the [asdf version manager](https://asdf-vm.com/) with the [asdf plugin for bashbot](https://github.com/mathew-fleisch/asdf-bashbot). It can also be used in [github actions](../.github/workflows/example-bashbot-github-action.yaml) as a simple text message stage, or as a [gating mechanism](../.github/workflows/example-bashbot-github-action-gate.yaml) to manually approve/reject a github action, or as a debugging breakpoint for troubleshooting a github action. These techniques should be applicable in any CI system.

```bash
# Download the default configuration files
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/.tool-versions -q -O ${PWD}/.tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml -q -O ${PWD}/config.yaml
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-env-file -q -O ${PWD}/.env

# Set a path to your local configuration file
export BASH_CONFIG_FILEPATH=${PWD}/config.yaml
# Set a path to your local configuration file
export BASH_CONFIG_FILEPATH=${PWD}/bashbot/config.yaml
# Set `Bot User OAuth Access Token` as SLACK_BOT_TOKEN environment variable
export SLACK_BOT_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
# Set `Bot App Access Token` as SLACK_APP_TOKEN environment variable
export SLACK_APP_TOKEN=xapp-xxxxxxxxx-xxxxxxx
# Set log-level (default: info)
export LOG_LEVEL=info
# Set log-format (default: text) [text|json]
export LOG_FORMAT=text

asdf plugin add bashbot
asdf install bashbot 2.0.2
asdf global bashbot 2.0.2

# To verify installation run version or help commands
bashbot version
# bashbot-darwin-amd64    v2.0.2

# See binary instructions after `bashbot verison`
```

### docker

```bash
# Set `Bot User OAuth Access Token` as SLACK_BOT_TOKEN environment variable
export SLACK_BOT_TOKEN=xoxb-xxxxxxxxx-xxxxxxx
# Set `Bot App Access Token` as SLACK_APP_TOKEN environment variable
export SLACK_APP_TOKEN=xapp-xxxxxxxxx-xxxxxxx

# Get the sample config.yaml
wget -O config.yaml https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml

# Pass environment variable and mount configuration yaml to run container
docker run -it --rm \
  --pull always \
  --name bashbot \
  -v ${PWD}/config.yaml:/bashbot/config.yaml \
  -e BASHBOT_CONFIG_FILEPATH="/bashbot/config.yaml" \
  -e SLACK_BOT_TOKEN \
  -e SLACK_APP_TOKEN \
  -e LOG_LEVEL="info" \
  -e LOG_FORMAT="text" \
  mathewfleisch/bashbot:latest
```

### kubernetes

To deploy bashbot with raw kubernetes manifests, see the rendered helm template for an example deployment with service account: [kubernetes-manifests-withsa.yaml](kubernetes/kubernetes-manifests-withsa.yaml)

Note: For bashbot to pick up changes to the configmap/values, a running pod must be deleted to trigger a restart.

```bash
# generated with
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/.tool-versions -q -O .tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml -q -O config.yaml

export NAMESPACE=bashbot
export BOTNAME=bashbot

helm template --debug \
  --version=${TARGET_VERSION} \
  --namespace ${NAMESPACE} \
  --set namespace=${NAMESPACE} \
  --set botname=${BOTNAME} \
  --set-file '\.tool-versions'=${PWD}/.tool-versions \
  --set-file 'config\.yaml'=${PWD}/config.yaml \
  --repo https://mathew-fleisch.github.io/bashbot \
  ${BOTNAME} bashbot > kubernetes-manifests-withsa.yaml
```

To deploy bashbot with raw kubernetes manifests, see the rendered helm template for an example deployment WITHOUT a service account: [kubernetes-manifests.yaml](kubernetes/kubernetes-manifests.yaml)

Note: For bashbot to pick up changes to the configmap/values, a running pod must be deleted to trigger a restart.

```bash
# generated with
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/.tool-versions -q -O .tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml -q -O config.yaml

export NAMESPACE=bashbot
export BOTNAME=bashbot

helm template --debug \
  --version=${TARGET_VERSION} \
  --namespace ${NAMESPACE} \
  --set namespace=${NAMESPACE} \
  --set botname=${BOTNAME} \
  --set serviceAccount.create=false \
  --set-file '\.tool-versions'=${PWD}/.tool-versions \
  --set-file 'config\.yaml'=${PWD}/config.yaml \
  --repo https://mathew-fleisch.github.io/bashbot \
  ${BOTNAME} bashbot > kubernetes/kubernetes-manifests.yaml
```

### helm

To install bashbot from the public helm repo, the slack app/bot tokens are required to be saved as kubernetes secret, and local versions of the config.yaml and .tool-versions file. To install the local helm chart from source, run the makefile target `make helm-install` after exporting the slack-app-token and slack-bot-token environment variables.

Note: For bashbot to pick up changes to the configmap/values, a running pod must be deleted to trigger a restart.

```bash
# Define your bashbot instance and namespace names
export BOTNAME=bashbot
export NAMESPACE=bashbot
export TARGET_VERSION=v2.0.2

# Get sample-env-file, sample-config.yaml, and .tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/.tool-versions -q -O ${PWD}/.tool-versions
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-config.yaml -q -O ${PWD}/config.yaml
wget https://raw.githubusercontent.com/mathew-fleisch/bashbot/main/sample-env-file -q -O ${PWD}/sample-env-file

# Add the public helm repo (one time)
helm repo add bashbot https://mathew-fleisch.github.io/bashbot

# Build kubernetes secrets from .env file (don't forget to fill the .env file with your secrets)
cat ${PWD}/sample-env-file | envsubst > ${PWD}/.env
kubectl --namespace ${NAMESPACE} create secret generic ${BOTNAME}-env \
  $(cat ${PWD}/.env | grep -vE '^#' | sed -e 's/export\ /--from-literal=/g' | tr '\n' ' ')

# Install a specific version of bashbot (remove --version for latest)
helm install --debug --wait \
  --version=${TARGET_VERSION} \
  --namespace ${NAMESPACE} \
  --set namespace=${NAMESPACE} \
  --set botname=${BOTNAME} \
  --set-file '\.tool-versions'=${PWD}/.tool-versions \
  --set-file 'config\.yaml'=${PWD}/config.yaml \
  --repo https://mathew-fleisch.github.io/bashbot \
  ${BOTNAME} bashbot
```

### KinD

For local development on bashbot, a kind cluster is used for testing.

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

### ArgoCD

Bashbot works great with ArgoCD because the values.yaml of the helm chart can override/set/customize any value, to customize for any deployment scenario. Note: updating the values in the config.yaml or .tool-versions section, only update a config-map, and the bashbot pod will need to be deleted, so the replacement pod will include the new configmap. [example-argocd-application.yam](example-argocd-application.yaml)

Note: For bashbot to pick up changes to the configmap/values, a running pod must be deleted to trigger a restart.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: bashbot
  namespace: argocd
spec:
  project: default
  source:
    chart: bashbot
    repoURL: https://mathew-fleisch.github.io/bashbot
    targetRevision: v2.0.2
    helm:
      releaseName: bashbot
      values: |
        botname: bashbot
        namespace: bots
        log_level: info
        log_format: text
        image:
          repository: mathewfleisch/bashbot
          pullPolicy: IfNotPresent
          command:
          - /bashbot/entrypoint.sh
        serviceAccount:
          create: true
        .tool-versions: |
          [contents of .tool-versions]
        config.yaml: |
          [contents of config.yaml]
  destination:
    server: "https://kubernetes.default.svc"
    namespace: bots
  syncPolicy:
    syncOptions:
        - CreateNamespace=true
    automated:
      prune: true
      selfHeal: true
```

---

## config.yaml

[sample-config.yaml](../sample-config.yaml)
The config.yaml file is defined as an array of yaml objects keyed by 'tools' and 'dependencies.' The dependencies section defines any resources that need to be downloaded or cloned from a repository before execution of commands. Each command is configured by a single object. For instance, the most simple command (hello world) returns a simple message when triggered, in the [ping example](ping).

```yaml
name: Ping/Pong
description: Return pong on pings
envvars: []
dependencies: []
help: "!bashbot ping"
trigger: ping
location: /bashbot/
command:
  - echo "pong"
parameters: []
log: true
ephemeral: false
response: text
permissions:
  - all
```

<img src="https://i.imgur.com/QyA6ECb.gif">

### Parameters and Arguments

Commands can take arguments, but require the values to either match a hard-coded list of values, value, from a sub-command, or a regular expression to be valid. Before each command is executed, a few environment variables are set, and can be leveraged in commands  `TRIGGERED_USER_NAME` `TRIGGERED_USER_ID` `TRIGGERED_CHANNEL_NAME` `TRIGGERED_CHANNEL_ID`

- name(string): title of command
- description(string): information about command
- envvars(array of strings): required environment variables to run command
- dependencies(array of strings): required pre-installed tools to run command
- help(string): message to users about how to execute
- trigger(string): single word to execute command
- location(string): absolute/relative path to set context
- command(array of strings): command to pass to bash -c
- parameters(array of objects):
  - name(string): title of command argument
  - allowed(array of strings): hard-coded list of allowed values
  - description(string): information about command argument
  - match(string): regular expression to match allowed values
  - source(array of strings): bash command to build list of allowed values
- log(boolean): whether to log command to log channel
- ephemeral(boolean): send response to user as ephemeral message
- response(text|code|file): send response in raw text, surrounded by markdown code tics, or as a text file
- permissions(array of strings): list of channel ids to restrict command to (all for unrestricted access) 

An example of a command that uses a hard-coded parameter to either run the `date` command or the `uptime` command on the host.

```yaml
name: Date or Uptime
description: Show the current time or uptime
envvars: []
dependencies: []
help: "!bashbot time"
trigger: time
location: /bashbot/
command:
  - "echo \"Date/time: $(${command})\""
parameters:
  - name: command
    allowed:
      - date
      - uptime
log: true
ephemeral: false
response: code
permissions:
  - all
```

An example of a command that uses `source` to build a list of allowed values is the [describe example](describe). Bashbot will first use the command line yaml parser yq to build a list of `trigger` values `yq e '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}`. A user would then only be able to "describe" each command that has a "trigger" value (i.e. `!bashbot describe describe` would show the following yaml)

```yaml
name: Describe Bashbot [command]
description: Show the yaml object for a specific command
envvars:
  - BASHBOT_CONFIG_FILEPATH
dependencies:
  - yq
help: "!bashbot describe [command]"
trigger: describe
location: /bashbot/
command:
  - yq e '.tools[] | select(.trigger=="${command}")' ${BASHBOT_CONFIG_FILEPATH}
parameters:
  - name: command
    allowed: []
    description: a command to describe ('bashbot list')
    source:
      - yq e '.tools[] | .trigger' ${BASHBOT_CONFIG_FILEPATH}
log: true
ephemeral: false
response: code
permissions:
  - all
```

<img src="https://i.imgur.com/bQZKRjX.gif">

An example of a command that has regular expression match parameter, is the [aqi example](aqi). This command takes one argument, a zip code and a regular expression is used to validate it. The [aqi script](aqi/aqi.sh) is used to curl the [Air Now API](https://docs.airnowapi.org/), and form a custom response with emojis.

```yaml
name: Air Quality Index
description: Get air quality index by zip code
envvars:
  - AIRQUALITY_API_KEY
dependencies:
  - curl
  - jq
help: "!bashbot aqi [zip-code]"
trigger: aqi
location: /bashbot/vendor/bashbot/examples/aqi
command:
  - "./aqi.sh ${zip}"
parameters:
  - name: zip
    allowed: []
    description: any zip code
    match: (^\d{5}$)|(^\d{9}$)|(^\d{5}-\d{4}$)
log: true
ephemeral: false
response: text
permissions:
  - all
```

<img src="https://i.imgur.com/GTgpdYf.png" />
