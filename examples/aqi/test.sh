#!/bin/bash
# shellcheck disable=SC2086,SC2181

set -o errexit
set -o pipefail
set -o nounset
# set -x # debug

cleanup() {
  echo "aqi test complete!"
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
      found_res=0
      for j in {3..1}; do
        bashbot_pod=$(kubectl -n ${ns} get pods -o jsonpath='{.items[0].metadata.name}')
        # Send `!bashbot info` via bashbot binary within bashbot pod
        kubectl --namespace ${ns} exec $bashbot_pod -- bash -c \
          'bashbot send-message --channel '${TESTING_CHANNEL}' --msg "!bashbot aqi 90210"'
        sleep 5
        last_log_line=$(kubectl -n ${ns} logs --tail 10 $bashbot_pod)
        # Tail the last line of the bashbot pod's log looking
        # for the string 'Air Quality Index' to prove the info script
        # is showing the correct value for whoami
        if [[ $last_log_line =~ "Air Quality Index" ]]; then
          echo "aqi test successful!"
          found_res=1

          kubectl --namespace ${ns} exec $bashbot_pod -- bash -c \
            'bashbot send-message --channel '${TESTING_CHANNEL}' --msg ":large_green_circle: aqi test successful!\nSaw \"Air Quality Index\" in bashbot logs"'
          exit 0
        fi
        echo "Bashbot aqi test failed. $j more attempts..."
        sleep 5
      done
      # Don't require aqi tests to pass for the whole test to pass
      [ $found_res -eq 1 ] && exit 0 || exit 1
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
    'bashbot send-message --channel '${TESTING_CHANNEL}' --msg ":red_circle: aqi test failed!"'
  exit 1
}

# Usage: ./test.sh [namespace] [deployment]
namespace=${1:-bashbot}
deploymentName=${2:-bashbot}

main $namespace $deploymentName
