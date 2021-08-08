#!/bin/bash
# shellcheck disable=SC2086,SC2154
set -eou pipefail

github_base="${github_base:-api.github.com}"
expected_variables="github_token github_org github_repo github_branch github_filename output_filename"
for expect in $expected_variables; do
  if [[ -z "${!expect}" ]]; then
    echo "Missing environment variable $expect"
    echo "Expected: $expected_variables"
    exit 1
  fi
done
echo "Downloading ${github_filename} from: https://${github_base}/repos/${github_org}/${github_repo}/contents/${github_filename}?ref=${github_branch}"
echo "To: ${output_filename}"
curl -H "Authorization: token $github_token" \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Content-Type: application/json" \
  -m 15 \
  -o ${output_filename} \
  -sL https://${github_base}/repos/${github_org}/${github_repo}/contents/${github_filename}?ref=${github_branch} 2>&1

