# Bashbot Example - Ping/Pong

In this example, a 'pong' response is returned to the user. This command is useful to verify Bashbot is running and as a template for hard-coded call-and-response message. For instance, common links, dates and/or milestones can be inserted in place of `echo \"pong\"` to share with your team.

## Bashbot configuration

This command is triggered by sending `bashbot ping` in a slack channel where Bashbot is also a member. This command requires no external script and simply returns a message 'pong' back to the user.

```json
{
  "name": "Ping/Pong",
  "description": "Return pong on pings",
  "help": "ping",
  "trigger": "ping",
  "location": "./",
  "setup": "echo \"\"",
  "command": "echo \"pong\"",
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "text",
  "permissions": ["all"]
}
```

