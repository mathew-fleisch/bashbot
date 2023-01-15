# Bashbot Example - Get Info

In this example, a bash script is executed to return information about the environment Bashbot is running in, and the channel/user name/id from which it was executed.

<img src="https://i.imgur.com/pIi7Pg2.gif">

## Bashbot configuration

This command is triggered by sending `bashbot info` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/version/get-info.sh` and requires no additional input to execute. It takes no arguments/parameters and returns `stdout` as a slack message, in the channel it was executed from.

```yaml
name: Get User/Channel Info
description: Get information about the user and channel command is being run from
envvars: []
dependencies: []
help: "!bashbot info"
trigger: info
location: /bashbot/vendor/bashbot/examples/info
command:
  - "./get-info.sh"
parameters: []
log: true
ephemeral: false
response: code
permissions:
  - all
```

## Bashbot script

Every Bashbot command, executed in slack, sets five environment variables describing "where," "when," and "who" executed the Bashbot command itself. The [get-info.sh](get-info.sh) script prints these environment variables and some other basic information about the host Bashbot is running from.