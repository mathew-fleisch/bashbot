#!/bin/bash
# shellcheck disable=SC2086,SC2181

set -o errexit
set -o pipefail
set -o nounset
# set -x # debug

cleanup() {
  echo "Regex test complete!"
}

trap cleanup EXIT
TESTING_CHANNEL=${TESTING_CHANNEL:-C034FNXS3FA}
main() {
  local ns=${1:-bashbot}
  local dn=${2:-bashbot}
  # Retry loop (20/$i with 3 second delay between loops)
  for i in {30..1}; do
    # Get the expected number of replicas for this deployment
    expectedReplicas=$(kubectl --namespace ${ns} -o jsonpath='{.status.replicas}' get deployment ${dn})
    # If the '.status.replicas' value is empty/not-set, set the default number of replicas to '1'
    [ -z "$expectedReplicas" ] && expectedReplicas=1
    # Get the number of replicas that are ready for this deployment
    readyReplicas=$(kubectl --namespace ${ns} -o jsonpath='{.status.readyReplicas}' get deployment ${dn})
    # If the .status.readyReplicas value is empty/not-set, set the default number of "ready" replicas to '0'
    [ -z "$readyReplicas" ] && readyReplicas=0

    # Test that the number of "ready" replicas match the number of expected replicas for this deployment
    test $readyReplicas -eq $expectedReplicas \
      && test 1 -eq 1
    if [ $? -eq 0 ]; then
      # echo "Bashbot deployment confirmed!"
      # kubectl --namespace ${ns} get deployments
      found_regex=0
      for j in {3..1}; do
        bashbot_pod=$(kubectl -n bashbot get pods -o jsonpath='{.items[0].metadata.name}')
        # Send `!bashbot curl https://jsonplaceholder.typicode.com/posts/1 | jq -r \‘.id\’`\
        # via bashbot binary within bashbot pod and expect an id of value 1 in the response
        kubectl --namespace ${ns} exec $bashbot_pod -- bash -c \
          'source .env && bashbot send-message --channel '${TESTING_CHANNEL}' --msg "!bashbot curl https://jsonplaceholder.typicode.com/posts/1 | jq -r ‘.id’"'
        sleep 5
        last_log_line=$(kubectl -n bashbot logs --tail 1 $bashbot_pod)
        # Tail the last line of the bashbot pod's log looking
        # for the string 'Bashbot is now connected to slack'
        if [[ $last_log_line =~ "1" ]]; then
          echo "regex test successful! Curl of jsonplaceholder.typicode.com api returned expected json object with .id=1, parsed by jq"
          found_regex=1

          kubectl --namespace ${ns} exec $bashbot_pod -- bash -c \
            'source .env && bashbot send-message --channel '${TESTING_CHANNEL}' --msg ":large_green_circle: regex test successful! Curl of jsonplaceholder.typicode.com api returned expected json object with .id=1, parsed by jq"'
          exit 0
        fi
        echo "Bashbot regex test failed. $j more attempts..."
        sleep 5
      done
      # Don't require regex tests to pass for the whole test to pass
      [ $found_regex -eq 1 ] && exit 0 || exit 1
    fi

    # Since the deployment was not ready, try again $i more times
    echo "Deployment not found or not ready. $i more attempts..."
    sleep 5
  done

  # The retry loop has exited without finding a stable deployment
  echo "Bashbot deployment failed :("
  # Display some debug information and fail test
  kubectl --namespace ${ns} get deployments
  kubectl --namespace ${ns} get pods -o wide

  kubectl --namespace ${ns} exec $bashbot_pod -- bash -c \
    'source .env && bashbot send-message --channel '${TESTING_CHANNEL}' --msg ":red_circle: regex test failed!"'
  exit 1
}

# Usage: ./test.sh [namespace] [deployment]
namespace=${1:-bashbot}
deploymentName=${2:-bashbot}

main $namespace $deploymentName
