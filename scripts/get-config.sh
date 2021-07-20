#!/bin/bash
# shellcheck disable=SC2086
set -eou pipefail

github_base="${github_base:-api.github.com}"
expected_variables="GIT_TOKEN github_org github_repo github_branch github_filename"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 1
  fi
done

curl -H "Authorization: token $GIT_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Content-Type: application/json" \
  -m 15 \
  -o ../config.json \
  -sL https://${github_base}/repos/${github_org}/${github_repo}/contents/${github_filename}?ref=${github_branch} 2>&1

