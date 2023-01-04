# Bashbot Example - Regex Command

In this example, a url is validated with a regular expression and curl is used to return the output as a message.

## Bashbot Configuration

This command is triggered by sending `bashbot curl [url]` in a slack channel where Bashbot is also a member. There is no external script for this command, and expects curl to already be installed.

```yaml
name: Curl Example
description: Pass a valid url to curl
envvars: []
dependencies:
  - curl
help: "!bashbot curl [url]"
trigger: curl
location: /bashbot/
command:
  - curl -s ${url} | jq -r ".body" | tr "\n" " "
parameters:
  - name: url
    allowed: []
    description: A valid url (Expecting json with key body)
    match: ^(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?$
log: true
ephemeral: false
response: code
permissions:
  - all
```

To match telephone numbers: `<tel:\+[0-9]+\|\+[0-9]+>`
