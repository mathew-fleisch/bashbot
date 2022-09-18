GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION?=$(shell git describe --abbrev=0 --tags)
LATEST_VERSION?=$(shell curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)
BINARY?=bin/bashbot
SRC_LOCATION?=cmd/bashbot/bashbot.go
LDFLAGS="-X main.Version=${VERSION}"
GO_BUILD=go build -ldflags=$(LDFLAGS)
BASHBOT_LOG_LEVEL?=info
BASHBOT_LOG_TYPE?=text
TESTING_CHANNEL?=C034FNXS3FA
NAMESPACE?=bashbot

.PHONY: help
help:
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
	@echo "| - make docker-build                                           |"
	@echo "|    + build and tag bashbot:local                              |"
	@echo "|                                                               |"
	@echo "| - make docker-run-local                                       |"
	@echo "|    + run an existing build of bashbot:local                   |"
	@echo "|                                                               |"
	@echo "| - make docker-run-local-bash                                  |"
	@echo "|    + run an exsting build of bashbot:local but                |"
	@echo "|      override the entrypoint with /bin/bash                   |"
	@echo "|                                                               |"
	@echo "| - make docker-run-upstream                                    |"
	@echo "|    + run the latest upstream build of bashbot                 |"
	@echo "|                                                               |"
	@echo "| - make docker-run-upstream-bash                               |"
	@echo "|    + run the latest upstream build of bashbot but             |"
	@echo "|      override the entrypoint with /bin/bash                   |"
	@echo "|                                                               |"
	@echo "| - make go-build                                               |"
	@echo "|    + build a go-binary for this host system-arch              |"
	@echo "|                                                               |"
	@echo "| - make go-build-cross-compile                                 |"
	@echo "|    + build go-binaries for linux/darwin amd64/arm64           |"
	@echo "|                                                               |"
	@echo "| - make go-clean                                               |"
	@echo "|    + delete any existing binaries at ./bin/*                  |"
	@echo "|                                                               |"
	@echo "| - make go-run                                                 |"
	@echo "|    + run the bashbot source code with go                      |"
	@echo "|                                                               |"
	@echo "| - make go-setup                                               |"
	@echo "|    + install go-dependencies                                  |"
	@echo "|                                                               |"
	@echo "| - make go-version                                             |"
	@echo "|    + run the bashbot source code with the --version           |"
	@echo "|      flag                                                     |"
	@echo "|                                                               |"
	@echo "| - make help                                                   |"
	@echo "|    + display this dialog                                      |"
	@echo "|                                                               |"
	@echo "| - make install-latest                                         |"
	@echo "|    + install the latest version of the bashbot binary to      |"
	@echo "|      /usr/local/bin/bashbot                                   |"
	@echo "|                                                               |"
	@echo "| - make pod-delete                                             |"
	@echo "|    + with an existing pod bashbot pod running, use kubectl    |"
	@echo "|      to delete it                                             |"
	@echo "|                                                               |"
	@echo "| - make pod-exec                                               |"
	@echo "|    + with an existing pod bashbot pod running, use kubectl    |"
	@echo "|      to exec into it                                          |"
	@echo "|                                                               |"
	@echo "| - make pod-get                                                |"
	@echo "|    + with an existing pod bashbot pod running, use kubectl    |"
	@echo "|      to get the pod name                                      |"
	@echo "|                                                               |"
	@echo "| - make pod-logs                                               |"
	@echo "|    + with an existing pod bashbot pod running, use kubectl    |"
	@echo "|      to display the logs of the pod                           |"
	@echo "|                                                               |"
	@echo "| - make test-docker                                            |"
	@echo "|    + use dockle to test the dockerfile for best practices     |"
	@echo "|                                                               |"
	@echo "| - make test-go                                                |"
	@echo "|    + run go coverage tests                                    |"
	@echo "|                                                               |"
	@echo "| - make test-kind                                              |"
	@echo "|    + run KinD tests                                           |"
	@echo "|                                                               |"
	@echo "| - make test-kind-cleanup                                      |"
	@echo "|    + delete any KinD cluster set up for bashbot               |"
	@echo "|                                                               |"
	@echo "| - make test-kind-helm-install                                 |"
	@echo "|    + install bashbot via helm into an existing KinD cluster   |"
	@echo "|      to /usr/local/bin/bashbot                                |"
	@echo "|                                                               |"
	@echo "| - make test-kind-setup                                        |"
	@echo "|    + setup a KinD cluster to test bashbot's helm chart        |"
	@echo "|                                                               |"
	@echo "| - make test-lint                                              |"
	@echo "|    + run golangci-lint                                        |"
	@echo "|                                                               |"
	@echo "| - make test-lint-actions                                      |"
	@echo "|    + run action-validator                                     |"
	@echo "|                                                               |"
	@echo "+---------------------------------------------------------------+"

.PHONY: docker-build
docker-build:
	docker build -t bashbot:local .

.PHONY: docker-run-local
docker-run-local:
	@docker run -it --rm \
		-v $(PWD)/config.json:/bashbot/config.json \
		-e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
		-v $(PWD)/.tool-versions:/bashbot/.tool-versions \
		-e SLACK_BOT_TOKEN=$(SLACK_BOT_TOKEN) \
		-e SLACK_APP_TOKEN=$(SLACK_APP_TOKEN) \
		-e LOG_LEVEL="$(BASHBOT_LOG_LEVEL)" \
		-e LOG_FORMAT="$(BASHBOT_LOG_TYPE)" \
		bashbot:local

.PHONY: docker-run-local-bash
docker-run-local-bash:
	@docker run -it --rm --entrypoint bash \
		-v $(PWD)/config.json:/bashbot/config.json \
		-e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
		-v $(PWD)/.tool-versions:/bashbot/.tool-versions \
		-e SLACK_BOT_TOKEN=$(SLACK_BOT_TOKEN) \
		-e SLACK_APP_TOKEN=$(SLACK_APP_TOKEN) \
		-e LOG_LEVEL="$(BASHBOT_LOG_LEVEL)" \
		-e LOG_FORMAT="$(BASHBOT_LOG_TYPE)" \
		bashbot:local

.PHONY: docker-run-upstream-bash
docker-run-upstream-bash:
	@docker run -it --rm --entrypoint bash \
		-v $(PWD)/config.json:/bashbot/config.json \
		-e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
		-e SLACK_BOT_TOKEN=$(SLACK_BOT_TOKEN) \
		-e SLACK_APP_TOKEN=$(SLACK_APP_TOKEN) \
		-e LOG_LEVEL="$(BASHBOT_LOG_LEVEL)" \
		-e LOG_FORMAT="$(BASHBOT_LOG_TYPE)" \
		mathewfleisch/bashbot:$(LATEST_VERSION)

.PHONY: docker-run-upstream
docker-run-upstream:
	@docker run -it --rm \
		-v $(PWD)/config.json:/bashbot/config.json \
		-e BASHBOT_CONFIG_FILEPATH="/bashbot/config.json" \
		-e SLACK_BOT_TOKEN=$(SLACK_BOT_TOKEN) \
		-e SLACK_APP_TOKEN=$(SLACK_APP_TOKEN) \
		-e LOG_LEVEL="$(BASHBOT_LOG_LEVEL)" \
		-e LOG_FORMAT="$(BASHBOT_LOG_TYPE)" \
		mathewfleisch/bashbot:$(LATEST_VERSION)

.PHONY: test-kind
test-kind: test-kind-setup test-kind-helm-install
	./helm/bashbot/test-deployment.sh
	./examples/ping/test.sh
	./examples/info/test.sh
	./examples/kubernetes/test.sh
	./helm/bashbot/test-complete.sh


.PHONY: test-kind-setup
test-kind-setup: docker-build
	kind create cluster || true
	kind load docker-image bashbot:local

.PHONY: test-kind-helm-install
test-kind-helm-install:
	helm upgrade bashbot helm/bashbot \
		--install \
		--namespace $(NAMESPACE) \
		--create-namespace \
		--set image.repository=bashbot \
		--set image.tag=local \
		--debug \
		--wait

.PHONY: test-kind-cleanup
test-kind-cleanup:
	helm --namespace $(NAMESPACE) delete bashbot || true
	kind delete cluster

.PHONY: pod-get
pod-get:
	@kubectl -n $(NAMESPACE) get pods | grep bashbot | cut -d' ' -f1

.PHONY: pod-logs
pod-logs:
	kubectl -n $(NAMESPACE) logs -f $(shell make pod-get) \
		| sed -e 's/\\n/\n/g'

.PHONY: pod-delete
pod-delete:
	kubectl -n $(NAMESPACE) delete pod $(shell make pod-get) --ignore-not-found=true

.PHONY: pod-exec
pod-exec:
	kubectl -n $(NAMESPACE) exec -it $(shell make pod-get) -- bash

.PHONY: test-lint-actions
test-lint-actions:
	find .github/workflows -type f \( -iname \*.yaml -o -iname \*.yml \) \
		| xargs -I {} action-validator --verbose {}

.PHONY: test-lint
test-lint:
	golangci-lint --verbose run

.PHONY: test-docker
test-docker:
	dockle bashbot:local

.PHONY: test-go
test-go:
	go test -cover -v ./...

.PHONY: go-setup
go-setup:
	go mod tidy
	go get github.com/slack-go/slack@master
	go get github.com/sirupsen/logrus
	go get -u golang.org/x/sys
	go mod vendor
	go install -v ./...


.PHONY: go-build-cross-compile
go-build-cross-compile: go-clean
	go mod tidy
	go mod vendor
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-linux-amd64  $(SRC_LOCATION)
	GOOS=linux   GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-linux-arm64  $(SRC_LOCATION)
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-darwin-amd64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-darwin-arm64 $(SRC_LOCATION)

.PHONY: go-build
go-build: go-clean
	go mod tidy
	go mod vendor
	CGO_ENABLED=0 $(GO_BUILD) -o $(BINARY)-$(GOOS)-$(GOARCH) $(SRC_LOCATION)

.PHONY: go-run
go-run:
	@go run $(SRC_LOCATION)

.PHONY: go-version
go-version:
	@go run $(SRC_LOCATION) --version

.PHONY: go-clean
go-clean:
	echo "Removing any existing go-binaries"
	rm -rf $(BINARY)*

.PHONY: install-latest
install-latest:
	wget -q -O /usr/local/bin/bashbot https://github.com/mathew-fleisch/bashbot/releases/download/$(LATEST_VERSION)/bashbot-$(GOOS)-$(GOARCH)
	chmod +x /usr/local/bin/bashbot
	bashbot --version
	@echo "Run 'bashbot --help' for more information"

.PHONY: gif
gif:
	@echo "Generating gif"
	@ffmpeg -i examples/$(example)/$(example).mov -r 10 -pix_fmt rgb24 examples/$(example)/$(example).gif

.PHONY: update-asdf-dependencies
update-asdf-dependencies:
	@curl -s -H "Accept: application/vnd.github.everest-preview+json" \
	    -H "Authorization: token $(GIT_TOKEN)" \
	    --request POST \
	    --data '{"event_type": "trigger-asdf-update"}' \
	    https://api.github.com/repos/mathew-fleisch/bashbot/dispatches
	@echo "Updating asdf dependencies via github-action: https://github.com/mathew-fleisch/bashbot/actions/workflows/update-asdf-versions.yaml"
