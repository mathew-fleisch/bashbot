#!/bin/bash
# shellcheck disable=SC2086
set -eou pipefail

config_file="${1:-}"
vendor_dir="${2:-}"

if [[ -z "$config_file" ]]; then
    echo "Must pass config file as first argument"
    exit 0
fi
if ! [[ -f "$config_file" ]]; then
    echo "Must pass config file as first argument."
    exit 0
fi
if [[ -z "$vendor_dir" ]]; then
    echo "Must pass vendor directory as second argument"
    exit 0
fi
if ! [[ -d "$vendor_dir" ]]; then
    echo "Must pass vendor directory as second argument"
    exit 0
fi
cd $vendor_dir && rm -rf *
echo "Installing Dependencies..."
jq -c '.dependencies[]' $config_file
for dep in $(jq -r '.dependencies[] | @base64' $config_file); do
  echo "---------------------"
    _jq() {
     echo ${dep} | base64 --decode | jq -r ${1}
    }

  this_name=$(_jq '.name')
  this_install=$(_jq '.install')

  echo "$this_name"
  echo "$this_install"
  eval $this_install
done
echo "Dependencies installed."
