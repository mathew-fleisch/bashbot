# Bashbot Example - Ping/Pong

In this example, a 'pong' response is returned to the user. This command is useful to verify Bashbot is running and as a template for hard-coded call-and-response message. For instance, common links, dates and/or milestones can be inserted in place of `echo \"pong\"` to share with your team.

<img src="https://i.imgur.com/QyA6ECb.gif">

## Bashbot configuration

This command is triggered by sending `bashbot ping` in a slack channel where Bashbot is also a member. This command requires no external script and simply returns a message 'pong' back to the user.

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
