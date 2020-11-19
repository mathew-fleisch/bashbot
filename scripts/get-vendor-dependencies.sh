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
dependencies=$(jq -r '.dependencies[].install' $config_file | sed 'N;s/\n/ \&\& /')
eval $dependencies
echo "Dependencies installed."