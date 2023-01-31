# Bashbot Example - Get File from (Private) Repo

In this example, a bash script is executed to download a file from a repository (that could be private) to the local file system. This command is useful for pulling fresh configuration without rebuilding bashbot.

## Bashbot configuration

This command is triggered by sending `bashbot get-file-from-repo` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/get-file-from-repo/get-file-from-repo.sh` and requires the following environment variables to be set: `github_token github_org github_repo github_branch github_filename output_filename`. This command requires [curl](https://curl.se/) to be installed on the host machine.

```json
{
  "name": "Update running configuration",
  "description": "Pulls a fresh configuration json file from github (could be private repo with GIT_TOKEN environment variable set)",
  "envvars": ["github_token", "github_org", "github_repo", "github_branch", "github_filename", "output_filename"],
  "dependencies": ["curl"],
  "help": "bashbot get-file-from-repo",
  "trigger": "get-file-from-repo",
  "location": "./examples/get-file-from-repo",
  "command": [
    "github_org=mathew-fleisch",
    "&& github_repo=bashbot",
    "&& github_filename=sample-config.yaml",
    "&& github_branch=main",
    "&& output_filename=${BASHBOT_CONFIG_FILEPATH}"
  ],
  "": "./get-file-from-repo.sh",
  "parameters": [],
  "log": false,
  "ephemeral": false,
  "response": "code",
  "permissions": ["all"]
}
```

## Bashbot script

This script expects environment variables to be set before executing and is designed to download a single file from a public or private repository. This command is useful to update the configuration json of a running instance of Bashbot without rebuilding or restarting. Each command that is executed in Bashbot will re-parse the json file at `$BASHBOT_CONFIG_FILEPATH`. Running this command will replace the running configuration json.

```bash
github_base="${github_base:-api.github.com}"
expected_variables="github_token github_org github_repo github_branch github_filename output_filename"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 1
  fi
done
echo "Downloading ${github_filename} from: https://api.github.com/repos/${github_org}/${github_repo}/contents/${github_filename}?ref=${github_branch}"
echo "To: ${output_filename}"
curl -H "Authorization: token $github_token" \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Content-Type: application/json" \
  -m 15 \
  -o ${output_filename} \
  -sL https://${github_base}/repos/${github_org}/${github_repo}/contents/${github_filename}?ref=${github_branch} 2>&1
```
