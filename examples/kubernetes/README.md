# Bashbot in Kubernetes

If Bashbot is installed with helm into a kubernetes cluster, with a service-account to make kube-api calls, it can be used as an SRE tool to quickly run (approved) troubleshooting kubectl commands to query status of the cluster that Bashbot is running in. Often times deleting specific pods is all a kubernetes cluster needs, to recover from a bad state (kubernetes version of "turning it off and on again"). Giving Bashbot the ability to carry out these actions (in private channels to restrict access within slack), can decrease time to resolution and provides a framework to codify runbooks as code.

## kubectl cluster-info

```yaml
name: Kubectl cluster-info
description: Return cluster-info by kubectl command
envvars: []
dependencies: []
help: "!bashbot k-info"
trigger: k-info
location: /bashbot/
command:
  - ". /usr/asdf/asdf.sh && kubectl cluster-info"
parameters: []
log: true
ephemeral: false
response: code
permissions:
  - all
```

## kubectl get pod

```yaml
name: kubectl get bashbot pod
description: Get the bashbot pod using kubectl
envvars:
  - BOTNAME
  - NAMESPACE
dependencies: []
help: "!bashbot k-get-pod"
trigger: k-get-pod
location: /bashbot/
command:
  - ". /usr/asdf/asdf.sh"
  - "&& kubectl --namespace ${NAMESPACE} get pods | grep -E \"NAME|${BOTNAME}\""
parameters: []
log: true
ephemeral: false
response: code
permissions:
  - all
```

## kubectl -n [namespace] delete pod [pod-name]

```yaml
name: kubectl -n [namespace] delete pod [podname]
description: Delete a pod
envvars: []
dependencies: []
help: "!bashbot k-delete-pod [namespace] [podname]"
trigger: k-delete-pod
location: /bashbot/
command:
  - ". /usr/asdf/asdf.sh"
  - "&& kubectl -n ${namespace} delete pod ${podname} --ignore-not-found=true"
parameters:
  - name: namespace
    allowed: []
    description: List all of the namespaces in the cluster
    source:
      - ". /usr/asdf/asdf.sh"
      - "&& kubectl get namespaces | grep -v NAME | awk '{print $1}'"
  - name: podname
    allowed: []
    description: List all of the pods in the cluster by name
    source:
      - ". /usr/asdf/asdf.sh"
      - "&& kubectl get pods -A | grep -v NAME | awk '{print $2}'"
log: true
ephemeral: false
response: code
permissions:
  - all
```

## kubectl -n [namespace] decribe pod [podname]

```yaml
name: kubectl -n [namespace] describe pod [podname]
description: Return pod information
envvars: []
dependencies: []
help: "!bashbot k-describe-pod [namespace] [podname]"
trigger: k-describe-pod
location: /bashbot/
command:
  - ". /usr/asdf/asdf.sh"
  - "&& kubectl -n ${namespace} describe pod ${podname}"
parameters:
  - name: namespace
    allowed: []
    description: List all of the namespaces in the cluster
    source:
      - ". /usr/asdf/asdf.sh"
      - "&& kubectl get namespaces | grep -v NAME | awk '{print $1}'"
  - name: podname
    allowed: []
    description: List all of the pods in the cluster by name
    source:
      - ". /usr/asdf/asdf.sh"
      - "&& kubectl get pods -A | grep -v NAME | awk '{print $2}'"
log: true
ephemeral: false
response: file
permissions:
  - all
```

## kubectl -n [namespace] logs -f [pod-name]

```yaml
name: kubectl -n [namespace] logs --tail 10 [podname]
description: Return last 10 lines of pod logs
envvars: []
dependencies: []
help: "!bashbot k-pod-logs [namespace] [podname]"
trigger: k-pod-logs
location: /bashbot/
command:
  - ". /usr/asdf/asdf.sh"
  - "&& kubectl -n ${namespace} logs --tail 10 ${podname}"
parameters:
  - name: namespace
    allowed: []
    description: List all of the namespaces in the cluster
    source:
      - ". /usr/asdf/asdf.sh"
      - "&& kubectl get namespaces | grep -v NAME | awk '{print $1}'"
  - name: podname
    allowed: []
    description: List all of the pods in the cluster by name
    source:
      - ". /usr/asdf/asdf.sh"
      - "&& kubectl get pods -A | grep -v NAME | awk '{print $2}'"
log: true
ephemeral: false
response: code
permissions:
  - all
```
