# Example Bashbot Commands

This directory contains a few sample Bashbot commands that can be applied to a configuration json by executing a helper script or by copy/pasting the json directly. Each example (directory) also contains a read-me file describing usage information. The [add-command.sh](add-command.sh) script expects the environment variable `BASH_CONFIG_FILEPATH` to be set and takes one argument (filepath of the command to add to the configuration json). The [remove-command.sh](remove-command.sh) script also expects the environment variable `BASH_CONFIG_FILEPATH` to be set and takes one argument (the value of the parameter `trigger` within the specific command).


```bash
$ export BASHBOT_CONFIG_FILEPATH=${PWD}/config.json

$ ./add-command.sh version/version.json                                           
# Trigger already exists: version

$ ./remove-command.sh version                
# version removed from /Users/user/bashbot/config.json
```


1. [Get Caller Information](info)
2. [Get Running Version](version)
3. [List Available Commands](list)
4. [Describe Command](describe)
5. [Get File From Repository](get-file-from-repo)





-------------------------------------------------------------------------

### config.json
[sample-config.json](sample-config.json)
The config.json file is defined as an array of json objects keyed by 'tools' and 'dependencies.' The dependencies section defines any resources that need to be downloaded or cloned from a repository before execution of commands. The following is a simplified example of a config.json file:

```json
{
  "admins": [{
    "trigger": "bashbot",
    "appName": "BashBot",
    "userIds": ["SLACK-USER-ID"],
    "privateChannelId": "SLACK-CHANNEL-ID",
    "logChannelId": "SLACK-CHANNEL-ID"
  }],
  "messages": [{
    "active": true,
    "name": "welcome",
    "text": "Witness the power of %s"
  },{
    "active": true,
    "name": "processing_command",
    "text": ":robot_face: Processing command..."
  },{
    "active": true,
    "name": "processing_raw_command",
    "text": ":smiling_imp: Processing raw command..."
  },{
    "active": true,
    "name": "command_not_found",
    "text": ":thinking_face: Command not found..."
  },{
    "active": true,
    "name": "incorrect_parameters",
    "text": ":face_with_monocle: Incorrect number of parameters"
  },{
    "active": true,
    "name": "invalid_parameter",
    "text": ":face_with_monocle: Invalid parameter value: %s"
  },{
    "active": true,
    "name": "ephemeral",
    "text": ":shushing_face: Message only shown to user who triggered it."
  },{
    "active": true,
    "name": "unauthorized",
    "text": ":skull_and_crossbones: You are not authorized to use this command in this channel.\nAllowed in: [%s]"
  }],
  "tools": [{
      "name": "List Commands",
      "description": "List all of the possible commands stored in bashbot",
      "help": "bashbot list-commands",
      "trigger": "list-commands",
      "location": "./",
      "setup": "echo \"\"",
      "command": "jq -r '.tools[] config.json | .trigger'",
      "parameters": [],
      "log": false,
      "ephemeral": false,
      "response": "code",
      "permissions": ["all"]
    }
  ],
  "dependencies": [
    {
      "name": "BashBot scripts Scripts",
      "install": "git clone https://github.com/mathew-fleisch/bashbot-scripts.git"
    }
  ]
}
```

Each object in the tools array defines the parameters of a single command.

```
name, description and help provide human readable information about the specific command
trigger:      unique alphanumeric word that represents the command
location:     absolute or relative path to dependency directory (use "./" for no dependency)
setup:        command that is run before the main command. (use "echo \"\"" as a default)
command:      bash command using ${parameter-name} to inject white-listed parameters or environment variables
parameters:   array of parameter objects. (more detail below)
log:          define whether the command should be logged in log channel
ephemeral:    define if the response should be shown to all, or just the user that triggered the command
response:     [code|text] code displays response in a code block, text displays response as raw text
permissions:  array of strings. private channel ids to restrict command access to
```

#### parameters continued (1 of 2):
There are a few ways to define parameters for a command, but the parameters passed to the bot MUST be white-listed. If the command can be triggered with no parameters, an empty array can be used as in the first example. If the command requires parameters, they can be a hard coded array of strings, or derived from another command. In this example, the hard-coded list of possible parameters is defined (random, question, answer). The `question` action essentially selects a random line in the `--questions-file` text file. The `answer` action does the same as questions, but with a greater-than sign appended to the string. Finally, the `random` action selects both, a random question and random answer from both linked text files.
```json
{
  "name": "Cards Against Humanity",
  "description": "Picks a random question and answer from a list.",
  "help": "bashbot cah [random|question|answer]",
  "trigger": "cah",
  "location": "./vendor/bashbot-scripts",
  "setup": "echo \"\"",
  "command": "./cardsAgainstHumanity.sh --action ${action} --questions-file ../against-humanity/questions.txt --answers-file ../against-humanity/answers.txt",
  "parameters": [{
    "name": "action",
    "allowed": ["random", "question", "answer"]
  }],
  "log": false,
  "ephemeral": false,
  "response": "text",
  "permissions": ["all"]
}
```