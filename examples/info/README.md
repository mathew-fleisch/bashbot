# Bashbot Example - Get Info

In this example, a bash script is executed to return information about the environment Bashbot is running in, and the channel/user name/id from which it was executed.

<img src="https://i.imgur.com/qC2ZnZ8.gif">

## Bashbot configuration

This command is triggered by sending `bashbot info` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/version/get-info.sh` and requires no additional input to execute. It takes no arguments/parameters and returns `stdout` as a slack message, in the channel it was executed from.

```json
{
  "name": "Get User/Channel Info",
  "description": "Get information about the user and channel command is being run from",
  "help": "info",
  "trigger": "info",
  "location": "./examples/info",
  "setup": "echo \"\"",
  "command": "./get-info.sh",
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```

## Bashbot script

Every Bashbot command, executed in slack, sets five environment variables describing "where," "when," and "who" executed the Bashbot command itself. The [get-info.sh](get-info.sh) script prints these environment variables and some other basic information about the host Bashbot is running from. The user name/id variables are useful in Bashbot commands to create some resource that ties back to a single user and message/tag those users in response messages. The channel name/id variables are useful to restrict certain Bashbot commands to a specific private channel to restrict execution access to only members of the private channel. The triggered-at variable is an epoch like timestamp and is useful for threading messages.

***Note: `exit 0` is used in success/failure states to ensure error messages are returned to slack. If `exit 1` is used for error states, a generic error message is returned to slack and `stderr` is suppressed.***

```bash
echo "Username[id]: ${TRIGGERED_USER_NAME}[${TRIGGERED_USER_ID}]"
echo " Channel[id]: ${TRIGGERED_CHANNEL_NAME}[${TRIGGERED_CHANNEL_ID}]"
echo "  Trigged at: ${TRIGGERED_AT}"
echo "------------------------------------------"
echo "        Date: $(date)"
echo "    uname -a: $(uname -a)"
echo "      uptime: $(uptime)"
echo "      whoami: $(whoami)"
```
