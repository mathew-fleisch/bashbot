```
---------------------------------------
 ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|
---------------------------------------
```
BashBot is a slack bot written in golang for infrastructure/devops teams. A socket connection to slack provides bashbot with a stream of text from each channel it is invited to, and uses regular expressions to trigger bash scripts. A [configuration file](sample-config.json) defines a list of commands that can be run in public and/or private channels. Restricting certain commands to private channels gives granular control, over which users can execute them. Bashbot allows infrastructure/devops teams to extend the tools and scripts they already use to manage their environments, into slack, that also acts as an execution log, and leverages slack's access controls.


--------------------------------------------------

## Installation and setup 

Bashbot can be run as a go binary or as a container and requires an .env file for secrets/environment-variables and a config.json saved in a _git repository_. The .env file will contain a slack token, a git token (for pulling private repositories), and the location of a config.json file. This _git repository_ should exist in your organization and should be devoted to your configuration of the bot. Bashbot will read from this repository constantly, making it easy to change the configuration without restarting the bot.



***Note***

Slack's permissions model has changed and the "[RTM](https://api.slack.com/rtm)" socket connection requires a "classic app" to be configured to get the correct type of token to run Bashbot. After logging into slack, visit [https://api.slack.com/apps?new_classic_app=1](https://api.slack.com/apps?new_classic_app=1) to set up a new "Bot User OAuth Access Token" and save the `xoxb-xxxxxxxxx-xxxxxxx` as the environment variable `SLACK_TOKEN` in a `.env` file at bashbot's root.

***Requirements***

- jq
- git
- golang


## Create .env and config.json

```bash
# Copy Sample Config
cp sample-config.json config.json

# Commit config.json to new repo (preferably private with token read access)

# Create .env file and fill in github_ variables pointing to custom config.json
touch .env
# Expected Format:
# export SLACK_TOKEN=xoxb-
# export GIT_TOKEN= <generate with permissions to READ from repo defined below>
# export github_org= <or username if you do not use an organizaion>
# export github_repo= <repo that will store the bot's config>
# export github_branch= <doesn't have to live on the main branch>
# export github_filename=path/to/config.json

# add secrets/tokens...
```

## Starting Manually

```bash
git clone git@github.com:mathew-fleisch/bashbot.git
# or
git clone https://github.com/mathew-fleisch/bashbot.git

cd bashbot
# Create/copy .env file and config.json to bashbot root

# Start local bashbot
./start.sh
```

## Starting via Docker

https://hub.docker.com/r/mathewfleisch/bashbot

```bash
# Create/copy .env file and config.json to wherever this next command runs:
docker run -v ${PWD}/config.json:/bashbot/config.json -v ${PWD}/.env:/bashbot/.env -it mathewfleisch/bashbot:v1.1.0
```


----------------------------------------------------------------


### .env file

```bash
export SLACK_TOKEN=xoxb-
export GIT_TOKEN=
export github_org=
export github_repo=
export github_branch=
export github_filename=path/to/config.json
```



### config.json
[sample-config.json](sample-config.json)
The config.json file is defined as an array of json objects keyed by 'tools' and 'dependencies.' The dependencies section defines any resources that need to be downloaded or cloned from a repository before execution of commands. The following is a simplified example of a config.json file:

```json
{
  "tools": [{
      "name": "List Commands",
      "description": "List all of the possible commands stored in bashbot",
      "help": "bashbot list-commands",
      "trigger": "list-commands",
      "location": "./",
      "setup": "echo \"\"",
      "command": "cat config.json | jq -r '.tools[] | .trigger' | sort",
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
      "install": "git clone https://$GITHUB_TOKEN@github.com/eaze/bashbot-scripts.git"
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

In this example, a user would type `bashbot list-commands` and that would then run the command `cat config.json | jq -r '.tools[] | .trigger' | sort` which takes no parameters and returns a code block of text from the response. 
```json
{
  "name": "List Commands",
  "description": "List all of the possible commands stored in bashbot",
  "help": "bashbot list-commands",
  "trigger": "list-commands",
  "location": "./",
  "setup": "echo \"\"",
  "command": "cat config.json | jq -r '.tools[] | .trigger' | sort",
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
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
#### parameters continued (2 of 2): 
In this example, a list of all 'trigger' values are extracted from the config.json and used as the parameter white-list. When the parameter list can be derived from output of another unix command, it can be "piped" in using the 'source' key. The command must execute without additional stdin input and consist of a newline separated list of values. The command jq is used to parse the json file of all 'trigger' values in a newline separated list.
```json
{
  "name": "Describe Command",
  "description": "Show the json object for a specific command",
  "help": "bashbot describe [command]",
  "trigger": "describe",
  "location": "./scripts",
  "setup": "echo \"\"",
  "command": "./describe-command.sh ../config.json ${command}",
  "parameters": [
    {
      "name": "command",
      "allowed": [],
      "description": "a command to describe ('bashbot list-commands')",
      "source": "cat ../config.json | jq -r '.tools[] | .trigger'"
    }],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```
