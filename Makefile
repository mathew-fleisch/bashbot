GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION?=$(shell make version)
LATEST_VERSION?=$(shell curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)
BINARY?=bin/bashbot
SRC_LOCATION?=main.go
LDFLAGS="-X github.com/mathew-fleisch/bashbot/cmd.Version=${VERSION}"
GO_BUILD=go build -ldflags=$(LDFLAGS)
# Public builds:  REGISTRY_NAME=mathewfleisch/bashbot or REGISTRY_NAME=ghcr.io/mathew-fleisch/bashbot
REGISTRY_NAME?=bashbot
# For latest tag: REGISTRY_TAG=latest"
REGISTRY_TAG?=local
NRUSER?=bb
BASHBOT_LOG_LEVEL?=info
BASHBOT_LOG_TYPE?=text
TESTING_CHANNEL?=C034FNXS3FA
ADMIN_CHANNEL?=GPFMM5MD2
NAMESPACE?=bashbot
BOTNAME?=bashbot
HELM_CONFIG_YAML?=$(PWD)/config.yaml
HELM_TOOL_VERSIONS?=$(PWD)/.tool-versions
HELM_ENV?=${PWD}/.env


##@ Misc stuff

.PHONY: help
help: ## this
	@echo "+---------------------------------------------------------------+"
	@echo "|   ____            _     ____        _                         |"
	@echo "|  |  _ \          | |   |  _ \      | |                        |"
	@echo "|  | |_) | __ _ ___| |__ | |_) | ___ | |_                       |"
	@echo "|  |  _ < / _' / __| '_ \|  _ < / _ \| __|                      |"
	@echo "|  | |_) | (_| \__ \ | | | |_) | (_) | |_                       |"
	@echo "|  |____/ \__,_|___/_| |_|____/ \___/ \__|                      |"
	@echo "|                                                               |"
	@echo "|  makefile targets                                             |"
	@echo "+---------------------------------------------------------------+"
	@echo "$(VERSION)"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: version
version: ## get the current helm chart version
	@yq e '.version' charts/bashbot/Chart.yaml

.PHONY: bump-patch
bump-patch: ## Bump-patch the semantic version of the helm chart using semver tool
	sed -i 's/'$(shell make version)'/v'$(shell semver bump patch $(shell make version))'/g' charts/bashbot/Chart.yaml

.PHONY: bump-minor
bump-minor: ## Bump-minor the semantic version of the helm chart using semver tool
	sed -i 's/'$(shell make version)'/v'$(shell semver bump minor $(shell make version))'/g' charts/bashbot/Chart.yaml

.PHONY: bump-major
bump-major: ## Bump-major the semantic version of the helm chart using semver tool
	sed -i 's/'$(shell make version)'/v'$(shell semver bump major $(shell make version))'/g' charts/bashbot/Chart.yaml

.PHONY: install-latest
install-latest: ## install the latest version of the bashbot binary to /usr/local/bin/bashbot with wget
	wget -q -O /usr/local/bin/bashbot https://github.com/mathew-fleisch/bashbot/releases/download/$(LATEST_VERSION)/bashbot-$(GOOS)-$(GOARCH)
	chmod +x /usr/local/bin/bashbot
ifeq ($(shell uname -s),Darwin)
	@echo "To add bashbot to an allowlist:"
	@echo "xattr -d com.apple.quarantine /usr/local/bin/bashbot"
endif
	bashbot version
	@echo "Run 'bashbot --help' for more information"

.PHONY: update-asdf-deps
update-asdf-deps: ## trigger github action to update asdf dependencies listed in .tool-versions (requires GIT_TOKEN)
	@curl -s -H "Accept: application/vnd.github.everest-preview+json" \
	    -H "Authorization: token $(GIT_TOKEN)" \
	    --request POST \
	    --data '{"event_type": "trigger-asdf-update"}' \
	    https://api.github.com/repos/mathew-fleisch/bashbot/dispatches
	@echo "Updating asdf dependencies via github-action: https://github.com/mathew-fleisch/bashbot/actions/workflows/update-asdf-versions.yaml"

.PHONY: gif
gif: ## Create a gif from a quicktime screen recording that has been exported to .mp4 from imovie
	@echo "Generating gif"
	@ffmpeg -i examples/$(example)/$(example).mp4 -r 10 -pix_fmt rgb24 examples/$(example)/$(example).gif

##@ Go stuff

.PHONY: go-build
go-build: go-clean go-setup ## build a go-binary for this host system-arch
	CGO_ENABLED=0 $(GO_BUILD) -o $(BINARY)-$(GOOS)-$(GOARCH) $(SRC_LOCATION)

.PHONY: go-clean
go-clean: ## delete any existing binaries at ./bin/*
	echo "Removing any existing go-binaries"
	rm -rf $(BINARY)*

.PHONY: go-setup
go-setup: ## install go-dependencies
	go mod tidy
	go get github.com/slack-go/slack@master
	go get github.com/sirupsen/logrus
	go get -u golang.org/x/sys
	go mod vendor
	go install -v ./...


.PHONY: go-cross-compile
go-cross-compile: go-clean go-setup ## build go-binaries for linux/darwin amd64/arm64
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-linux-amd64  $(SRC_LOCATION)
	GOOS=linux   GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-linux-arm64  $(SRC_LOCATION)
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-darwin-amd64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-darwin-arm64 $(SRC_LOCATION)

.PHONY: go-run
go-run: ## run the bashbot source code with go
	@go run $(SRC_LOCATION)

.PHONY: go-version
go-version: ## run the bashbot source code with the version argument
	@go run $(SRC_LOCATION) version


##@ Docker stuff


.PHONY: docker-build
docker-build: ## build and tag bashbot:local

	docker build --build-arg NRUSER=$(NRUSER) -t $(REGISTRY_NAME):$(REGISTRY_TAG) .

.PHONY: docker-run
docker-run: ## run an existing build of $(REGISTRY_NAME):$(REGISTRY_TAG)
	@echo "Public Builds:  REGISTRY_NAME=mathewfleisch/bashbot or REGISTRY_NAME=ghcr.io/mathew-fleisch/bashbot"
	@echo "for latest tag: REGISTRY_TAG=latest"
	docker run -it --rm \
		-v $(HELM_CONFIG_YAML):/bashbot/config.yaml \
		-e BASHBOT_CONFIG_FILEPATH="/bashbot/config.yaml" \
		-v $(HELM_TOOL_VERSIONS):/bashbot/.tool-versions \
		-e GIT_TOKEN \
		-e AIRQUALITY_API_KEY \
		-e SLACK_BOT_TOKEN \
		-e SLACK_APP_TOKEN \
		-e LOG_LEVEL="$(BASHBOT_LOG_LEVEL)" \
		-e LOG_FORMAT="$(BASHBOT_LOG_TYPE)" \
		$(REGISTRY_NAME):$(REGISTRY_TAG)

.PHONY: docker-run-bash
docker-run-bash: ## run an exsting build of $(REGISTRY_NAME):$(REGISTRY_TAG) but override the entrypoint with /bin/bash
	@echo "Public Builds:  REGISTRY_NAME=mathewfleisch/bashbot or REGISTRY_NAME=ghcr.io/mathew-fleisch/bashbot"
	@echo "for latest tag: REGISTRY_TAG=latest"
	docker run -it --rm --entrypoint bash \
		-v $(HELM_CONFIG_YAML):/bashbot/config.yaml \
		-e BASHBOT_CONFIG_FILEPATH="/bashbot/config.yaml" \
		-v $(HELM_TOOL_VERSIONS):/bashbot/.tool-versions \
		-e GIT_TOKEN \
		-e AIRQUALITY_API_KEY \
		-e SLACK_BOT_TOKEN \
		-e SLACK_APP_TOKEN \
		-e LOG_LEVEL="$(BASHBOT_LOG_LEVEL)" \
		-e LOG_FORMAT="$(BASHBOT_LOG_TYPE)" \
		$(REGISTRY_NAME):$(REGISTRY_TAG)



##@ Kubernetes stuff


.PHONY: test-kind
test-kind: kind-setup helm-install test-run ## run KinD tests

test-run: ## run tests designed for bashbot running in kubernetes
	@echo "Testing: $(NAMESPACE) $(BOTNAME)"
	./charts/bashbot/test-deployment.sh $(NAMESPACE) $(BOTNAME)
	./examples/ping/test.sh $(NAMESPACE) $(BOTNAME)
	./examples/aqi/test.sh $(NAMESPACE) $(BOTNAME) || true
	./examples/asdf/test.sh $(NAMESPACE) $(BOTNAME)
	./examples/info/test.sh $(NAMESPACE) $(BOTNAME)
	./examples/regex/test.sh $(NAMESPACE) $(BOTNAME)
	./examples/kubernetes/test.sh $(NAMESPACE) $(BOTNAME)
	./charts/bashbot/test-complete.sh $(NAMESPACE) $(BOTNAME)

.PHONY: kind-setup
kind-setup: docker-build ## setup a KinD cluster to test bashbot's helm chart
	kind create cluster || true
	kind load docker-image $(REGISTRY_NAME):$(REGISTRY_TAG)

.PHONY: kind-cleanup
kind-cleanup: ## delete any KinD cluster set up for bashbot
	kind delete cluster


.PHONY: helm-install
helm-install: helm-uninstall ## install bashbot via helm into an existing KinD cluster to /usr/local/bin/bashbot
	kubectl create namespace $(NAMESPACE) || true
	@echo "Creating kubernetes secrets from $(HELM_ENV)"
	@echo "kubectl --namespace $(NAMESPACE) get secret $(BOTNAME)-env"
	@kubectl --namespace $(NAMESPACE) create secret generic $(BOTNAME)-env \
		$(shell cat $(HELM_ENV) | grep -vE '^#' | sed -e 's/export\ /--from-literal=/g' | tr '\n' ' ');
	helm upgrade $(BOTNAME) charts/bashbot \
		--install \
		--timeout 2m0s \
		--namespace $(NAMESPACE) \
		--set namespace=$(NAMESPACE) \
		--set botname=$(BOTNAME) \
		--set image.repository=$(REGISTRY_NAME) \
		--set image.tag=$(REGISTRY_TAG) \
		--set log_level=$(BASHBOT_LOG_LEVEL) \
		--set log_format=$(BASHBOT_LOG_TYPE) \
		--set-file '\.tool-versions'=$(HELM_TOOL_VERSIONS) \
		--set-file 'config\.yaml'=$(HELM_CONFIG_YAML) \
		--debug \
		--wait
	sleep 3
	timeout 30s make pod-logs 2>/dev/null || sleep 30

.PHONY: helm-uninstall
helm-uninstall: ## uninstall bashbot via helm/kubectl from an existing cluster
	helm uninstall $(BOTNAME) --namespace $(NAMESPACE) 2>/dev/null || true
	kubectl --namespace $(NAMESPACE) delete secret $(BOTNAME)-env --ignore-not-found=true
	kubectl delete clusterrolebinding $(BOTNAME) --ignore-not-found=true
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true

.PHONY: pod-get
pod-get: ## with an existing pod bashbot pod running, use kubectl to get the pod name
	@kubectl --namespace $(NAMESPACE) get pods \
		--template '{{range .items}}{{.metadata.name}}{{end}}' \
		--selector=app=$(BOTNAME)

.PHONY: pod-logs
pod-logs: ## with an existing pod bashbot pod running, use kubectl to display the logs of the pod
	kubectl -n $(NAMESPACE) logs -f $(shell make pod-get) \
		| sed -e 's/\\*\\n/\n/g'
	
.PHONY: pod-logs-json
pod-logs-json: ## with an existing pod bashbot pod running, use kubectl to display the json logs of the pod and pipe to jq
	kubectl -n $(NAMESPACE) logs -f $(shell make pod-get) \
		| jq -Rr '. as $$line | try (fromjson | .) catch $$line'

.PHONY: pod-delete
pod-delete: ## with an existing pod bashbot pod running, use kubectl to delete it
	kubectl -n $(NAMESPACE) delete pod $(shell make pod-get) --ignore-not-found=true

.PHONY: pod-exec
pod-exec: ## with an existing pod bashbot pod running, use kubectl to exec into it 
	kubectl -n $(NAMESPACE) exec -it $(shell make pod-get) -- bash

.PHONY: pod-exec-test
pod-exec-test: ## with an existing pod bashbot pod running, use kubectl to exec into it and run the test-suite
	kubectl -n $(NAMESPACE) exec  $(shell make pod-get) -- \
		bash -c '. /usr/asdf/asdf.sh && make test-run'


##@ Linters and Tests


.PHONY: test-lint-actions
test-lint-actions: ## lint github actions with action-validator
	find .github/workflows -type f \( -iname \*.yaml -o -iname \*.yml \) \
		| xargs -I {} action-validator --verbose {}

.PHONY: test-lint
test-lint: ## lint go source with golangci-lint
	golangci-lint --skip-dirs-use-default --verbose run || true

.PHONY: test-docker
test-docker: ## use dockle to test the dockerfile for best practices
	export DOCKER_CONTENT_TRUST=1 \
		&& make docker-build \
		&& dockle $(REGISTRY_NAME):$(REGISTRY_TAG)

# go test -cover -v ./...
.PHONY: test-go
test-go: ## run go coverage tests
	@echo "no tests..."
