#!/bin/bash
# shellcheck disable=SC2086
set -eou pipefail

expected_variables="GIT_TOKEN github_org github_repo github_branch github_filename"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    exit 1
  fi
done

curl -H "Authorization: token $GIT_TOKEN" \
  -H 'Accept: application/vnd.github.v3.raw' \
  -o ../config.json \
  -sL https://api.github.com/repos/${github_org}/${github_repo}/contents/${github_filename}?ref=${github_branch} 2>&1

