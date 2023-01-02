# Examples

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

<img src="https://i.imgur.com/uLrCYTf.gif">

## Parameters

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

---

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
