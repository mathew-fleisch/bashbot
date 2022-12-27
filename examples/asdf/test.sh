#!/bin/bash
# shellcheck disable=SC2086,SC2181

set -o errexit
set -o pipefail
set -o nounset
# set -x # debug

cleanup() {
  echo "asdf test complete!"
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
      found_asdf=0
      for j in {3..1}; do
        bashbot_pod=$(kubectl -n ${ns} get pods -o jsonpath='{.items[0].metadata.name}')
        asdf_found=$(kubectl --namespace ${ns} exec $bashbot_pod -- bash -c '. /usr/asdf/asdf.sh && command -v asdf')
        # Tail the last line of the bashbot pod's log looking
        # for the string 'Bashbot is now connected to slack'
        if [ -n "$asdf_found" ]; then
          echo "asdf found!"
          found_asdf=1

          kubectl --namespace ${ns} exec $bashbot_pod -- bash -c \
            'bashbot send-message --channel '${TESTING_CHANNEL}' --msg ":large_green_circle: asdf installed successfully!"'
          exit 0
        fi
        echo "Bashbot dependency test failed (asdf). $j more attempts..."
        sleep 5
      done
      [ $found_asdf -eq 1 ] && exit 0 || exit 1
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
    'bashbot send-message --channel '${TESTING_CHANNEL}' --msg ":red_circle: dependency test (asdf) failed!"'
  exit 1
}

# Usage: ./test.sh [namespace] [deployment]
namespace=${1:-bashbot}
deploymentName=${2:-bashbot}

main $namespace $deploymentName
