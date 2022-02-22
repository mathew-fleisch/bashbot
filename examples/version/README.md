# Bashbot Example - Get Version

In this example, a bash script is executed, attempting three methods of returning the currently running version of Bashbot.

<img src="https://i.imgur.com/ZQmH672.gif">

## Bashbot configuration

This command is triggered by sending `bashbot version` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/version/get-version.sh` and requires no additional input to execute. It takes no arguments/parameters and returns `stdout` as a slack message, in the channel it was executed from.

```json
{
  "name": "Get Bashbot Version",
  "description": "Displays the currently running version of Bashbot",
  "help": "bashbot version",
  "trigger": "version",
  "location": "./examples/version",
  "command": ["./get-version.sh"],
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```

## Bashbot script

The [get-version.sh](get-version.sh) first checks to see if Bashbot is already available via `$PATH` and prints that version, if possible. Next, if the binary is being run from the `./bin` directory, this script will grab the correct os/architecture from `uname` and attempt to print that version. If both of those methods fail, a [Makefile](../../Makefile) target (`go-version`) is executed, provided the Makefile and go source exists. If all of those methods are unsuccessful, an error message is returned to slack that the version could not be determined.

***Note: `exit 0` is used in success/failure states to ensure error messages are returned to slack. If `exit 1` is used for error states, a generic error message is returned to slack and `stderr` is suppressed.***

```bash
# First check if bashbot is installed and pull that version
if command -v bashbot > /dev/null; then 
    command -v bashbot
    bashbot --version
    exit 0
fi

# Next check if ./bin/bashbot-${os}-${arch} exists and pull that version
arch=amd64
[ "$(uname -m)" == "aarch64" ] && arch=arm64
os=$(uname | tr '[:upper:]' '[:lower:]')
if [ -f "./bin/bashbot-${os}-${arch}" ]; then
    echo "./bin/bashbot-${os}-${arch} --version"
    ./bin/bashbot-${os}-${arch} --version
    exit 0
fi

# Finally, check if bashbot go source exists and pull that version
go_filename=cmd/bashbot/bashbot.go
if [[ -f "./${go_filename}" ]] && [[ -f "./Makefile" ]]; then
    make go-setup
    make go-version
    exit 0
fi

echo "Could not determine the current version of bashbot"
exit 0
```
