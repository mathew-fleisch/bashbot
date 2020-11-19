#!/bin/bash

config_file="${1:-}"
trigger="${2:-}"

if [[ -z $config_file ]]; then
  echo "Missing config file"
  exit 0
fi
if ! [[ -f "$config_file" ]]; then
  echo "File not found..."
  exit 0
fi
if [[ -z $trigger ]]; then
  echo "Missing command to describe"
  exit 0
fi


this_config="$(cat $config_file | jq -cr '.tools[] | select(.trigger == "'${trigger}'")')"
new_params=""
while read param; do
  if [[ -n $(echo $param | grep source) ]]; then
    this_name=$(echo "$param" | jq -r '.name')
    this_allowed=$(echo "$param" | jq -r '.allowed')
    this_description=$(echo "$param" | jq -r '.description')
    this_source=$(echo "$param" | jq -r '.source')
    this_source_result="$($this_source)"
    this_allowed="$(echo "$this_source_result" | tr '\n' ',' | sed -e 's/,/","/g' | sed -e 's/^\(.*\)$/"\1"/g' | sed -e 's/,""//g')"
    new_params="$new_params,{\"name\":\"$this_name\",\"allowed\":[$this_allowed],\"description\":\"$this_description\",\"source\":\"$this_source\",}"
  else
    new_params="$new_params,$param"
  fi
done <<< "$(echo "$this_config" | jq -c '.parameters[]')"
new_params="$(echo "$new_params" | sed -e 's/^,//g')"
echo "$this_config" | jq '.parameters |= ['"$new_params"']'
