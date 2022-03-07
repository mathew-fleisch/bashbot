# Bashbot in Kubernetes

If Bashbot is installed with helm into a kubernetes cluster, with a service-account to make kube-api calls, it can be used as an SRE tool to quickly run (approved) troubleshooting kubectl commands to query status of the cluster that Bashbot is running in. Often times deleting specific pods is all a kubernetes cluster needs, to recover from a bad state (kubernetes version of "turning it off and on again"). Giving Bashbot the ability to carry out these actions (in private channels to restrict access within slack), can decrease time to resolution and provides a framework to codify runbooks as code.

