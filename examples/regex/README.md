# Bashbot Example - Regex Command

In this example, a url is validated with a regular expression and curl is used to return the output as a message.

## Bashbot Configuration

This command is triggered by sending `bashbot curl [url]` in a slack channel where Bashbot is also a member. There is no external script for this command, and expects curl to already be installed. Passing a valid url to this curl command will return the http response of that url, using the host system's built-in curl command. Note: Backslashes `\` must be escaped with another backslash. Example: matching a word `\w` would require passing `\\w`

```json
{
  "name": "Curl Example",
  "description": "Pass a valid url to curl",
  "help": "bashbot curl [url]",
  "trigger": "curl",
  "location": "./",
  "command": ["curl -s ${url}"],
  "parameters": [
    {
      "name": "url",
      "allowed": [],
      "description": "A valid url",
      "match": "(http|ftp|https)://([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?"
    }
  ],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```

To match telephone numbers:

```json
    {
      "name": "phone_number",
      "description": "Telephone number (format: +123456789)",
      "match": "<tel:\\+[0-9]+\\|\\+[0-9]+>"
    }
```
