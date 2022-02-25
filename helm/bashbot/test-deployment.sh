#!/bin/bash
# shellcheck disable=SC2086,SC2181

set -o errexit
set -o pipefail
set -o nounset
# set -x # debug

cleanup() {
  echo "Helm test complete!"
}

trap cleanup EXIT

main() {
  local ns=${1:-bashbot}
  local dn=${2:-bashbot}
  # Retry loop (20/$i with 3 second delay between loops)
  for i in {20..1}; do
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
      echo "Bashbot deployment confirmed!"
      kubectl --namespace ${ns} get deployments
      exit 0
    fi

    # Since the deployment was not ready, try again $i more times
    echo "Deployment not found or not ready. $i more attempts..."
    sleep 3
  done

  # The retry loop has exited without finding a stable deployment
  echo "Bashbot deployment failed :("
  # Display some debug information and fail test
  kubectl --namespace ${ns} get deployments
  kubectl --namespace ${ns} get pods -o wide
  exit 1
}

# Usage: ./test-deployment.sh [namespace] [deployment]
namespace=${1:-bashbot}
deploymentName=${2:-bashbot}

main $namespace $deploymentName
