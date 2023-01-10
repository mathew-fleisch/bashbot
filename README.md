# Bashbot

[![Release](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml) |
[Docker Hub](https://hub.docker.com/r/mathewfleisch/bashbot/tags?page=1&ordering=last_updated) | [ghcr](https://github.com/mathew-fleisch/bashbot/pkgs/container/bashbot)

BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to determine when to trigger bash commands. A [configuration file](sample-config.yaml) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls.

See the [examples](examples) directory for more information about deploying, configuring and customizing Bashbot for your team. In this example, a user triggers a Jenkins job using Bashbot and another instance of Bashbot is deployed in a Jenkins job as a gating mechanism. The configuration for the secondary Bashbot could get info about the Jenkins job/host and provides controls to manually decide if the job should pass or fail, at a certain stage in the build. This method of deploying Bashbot gives basic Jenkins controls (trigger, pass, fail) to users in an organization, without giving them access to Jenkins itself. Bashbot commands can be restricted to private channels to limit access within slack.

<img src="https://i.imgur.com/P6IL10y.gif" />

---

## Slack tokens

Bashbot uses the Slack API's "[Socket Mode](https://api.slack.com/apis/connections/socket)" to connect to the slack servers over a socket connection and uses ***no webhooks/ingress*** to trigger commands. Bashbot "subscribes" to message and emoji events and determins if a command should be executed and what command should be executed by parsing the data in each event. To run Bashbot, you must have a "[Bot User OAuth Token](https://api.slack.com/authentication/token-types#bot)" and an "[App-Level Token](https://api.slack.com/authentication/token-types#app)".

- Click "Create New App" from the [Slack Apps](https://api.slack.com/apps) page and follow the "From Scratch" prompts to give your instance of bashbot a unique name and a workspace for it to be installed in
- The "Basic Information" page gives controls to set a profile picture toward the bottom (make sure to save any changes)
- Enable "Socket Mode" from the "Socket Mode" page and add the default scopes `conversations.write` and note the "[App-Level Token](https://api.slack.com/authentication/token-types#app)" that is generated to save in the .env file as `SLACK_APP_TOKEN`
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

---

## Installation, setup and configuration

[Example deployments and commands](examples)

Bashbot can be run as a go binary or as a container and requires a slack-bot-token, slack-app-token and a config.yaml. The go binary takes flags to set the slack-bot-token, slack-app-token and path to the config.yaml file and the container uses environment variables to trigger a go binary by [entrypoint.sh](entrypoint.sh). The [Examples](examples) directory of this repository, has many deployment examples, commands used in automated tests, and can illustrate how bashbot can be used to trigger automation, or as automation itself, by leveraging secrets from the host running bashbot. For instance, one command might use an api token to curl a third-party api, and return the response back to the slack user, once triggered ([aqi example](examples/aqi)). If deployed in a kubernetes cluster, with a service-account to access the kube-api, bashbot can allow slack users to execute (hopefully curated) kubectl/helm commands, for an SRE/Ops focused deployment (get/describe/delete pods/deployments/secrets etc).

***Steps To Prove It's Working***

- Invite the BashBot into a channel by typing `@BashBot`
- Slackbot should respond with the message: `OK! I’ve invited @BashBot to this channel.`
- Now type `!bashbot help`
- If all is configured correctly, you should see BashBot respond immediately with `Processing command...` and momentarily post a full list of commands that are defined in config.yaml

---

## Automation

Included in this repository, one github action is executed on commits to pull requests, and another is executed on merges to the main branch:

1. On [pull requests to the main branch](.github/workflows/pr.yaml) [![pr](https://github.com/mathew-fleisch/bashbot/actions/workflows/pr.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/pr.yaml), four jobs are run on every commit:
    - linting and unit tests are run under the `unit_tests` job
    - a container is built and scanned by the [anchore](https://anchore.com/) container scanning tool
    - the golang code is analyzed by [codeql](https://codeql.github.com/) SAST tool
    - a container is built and deployed in a kind cluster, to run automated tests, to maintain/verify basic functionality (see [Makefile](Makefile) target `make help` for more information)
2. The [![release](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/release.yaml) action will:
    - cross compile go-binaries for linux/amd64, linux/arm64, darwin/amd64, and darwin/arm64
    - package and release a helm chart with the chart-releaser action
    - add the go-binaries as release artifacts
    - use the buildx docker plugin to build and push a container for amd64/arm64 to docker hub and ghcr.
    - use asdf to install the new version of bashbot
    - notify the bashbot test slack workspace
3. The [![update-asdf-versions](https://github.com/mathew-fleisch/bashbot/actions/workflows/update-asdf-versions.yaml/badge.svg)](https://github.com/mathew-fleisch/bashbot/actions/workflows/update-asdf-versions.yaml) action will check for new versions of dependencies installed with the [asdf version manager](https://asdf-vm.com/), found in the [.tool-versions](.tool-versions) file.

### Makefile

run `make help` for a full list of targets.

Note: [yq](https://github.com/mikefarah/yq) is a dependency of running many makefile targets and can be installed with the binary: `wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -q -O /usr/local/bin/yq && chmod +x /usr/local/bin/yq`

```text
+---------------------------------------------------------------+
|   ____            _     ____        _                         |
|  |  _ \          | |   |  _ \      | |                        |
|  | |_) | __ _ ___| |__ | |_) | ___ | |_                       |
|  |  _ < / _' / __| '_ \|  _ < / _ \| __|                      |
|  | |_) | (_| \__ \ | | | |_) | (_) | |_                       |
|  |____/ \__,_|___/_| |_|____/ \___/ \__|                      |
|                                                               |
|  makefile targets                                             |
+---------------------------------------------------------------+
v2.0.2

Usage:
  make <target>

Go stuff
  go-build            build a go-binary for this host system-arch
  go-clean            delete any existing binaries at ./bin/*
  go-setup            install go-dependencies
  go-cross-compile    build go-binaries for linux/darwin amd64/arm64
  go-run              run the bashbot source code with go
  go-version          run the bashbot source code with the version argument

Docker stuff
  docker-build        build and tag bashbot:local
  docker-run          run an existing build of bashbot:local
  docker-run-bash     run an exsting build of bashbot:local but override the entrypoint with /bin/bash
  docker-run-up-bash  run the latest upstream build of bashbot but override the entrypoint with /bin/bash
  docker-run-up       run the latest upstream build of bashbot

Kubernetes stuff
  test-kind           run KinD tests
  test-run            run tests designed for bashbot running in kubernetes
  kind-setup          setup a KinD cluster to test bashbot's helm chart
  kind-cleanup        delete any KinD cluster set up for bashbot
  version             get the current helm chart version
  helm-bump-patch     Bump-patch the semantic version of the helm chart using semver tool
  helm-bump-minor     Bump-minor the semantic version of the helm chart using semver tool
  helm-bump-major     Bump-major the semantic version of the helm chart using semver tool
  helm-install        install bashbot via helm into an existing KinD cluster to /usr/local/bin/bashbot
  helm-uninstall      uninstall bashbot via helm/kubectl from an existing cluster
  pod-get             with an existing pod bashbot pod running, use kubectl to get the pod name
  pod-logs            with an existing pod bashbot pod running, use kubectl to display the logs of the pod
  pod-logs-json       with an existing pod bashbot pod running, use kubectl to display the json logs of the pod and pipe to jq
  pod-delete          with an existing pod bashbot pod running, use kubectl to delete it
  pod-exec            with an existing pod bashbot pod running, use kubectl to exec into it 
  pod-exec-test       with an existing pod bashbot pod running, use kubectl to exec into it and run the test-suite

Linters and Tests
  test-lint-actions   lint github actions with action-validator
  test-lint           lint go source with golangci-lint
  test-docker         use dockle to test the dockerfile for best practices
  test-go             run go coverage tests

Other stuff
  help                this
  install-latest      install the latest version of the bashbot binary to /usr/local/bin/bashbot with wget
  update-asdf-deps    trigger github action to update asdf dependencies listed in .tool-versions (requires GIT_TOKEN)
```

Any of these environment variables can be overridden by exporting a new value before running any makefile target.

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

