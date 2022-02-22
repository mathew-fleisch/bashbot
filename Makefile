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

.PHONY: go-setup
go-setup:
	go mod tidy
	go mod vendor
	go install -v ./...
	go get github.com/slack-go/slack@master
	go get github.com/sirupsen/logrus

.PHONY: go-build-cross-compile
go-build-cross-compile:
	rm -rf $(BINARY)*
	go mod tidy
	go mod vendor
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-linux-amd64 $(SRC_LOCATION)
	GOOS=linux   GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-linux-arm64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-darwin-amd64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-darwin-arm64 $(SRC_LOCATION)

.PHONY: go-build
go-build:
	rm -rf $(BINARY)-$(GOOS)-$(GOARCH)
	go mod tidy
	go mod vendor
	CGO_ENABLED=0 $(GO_BUILD) -o $(BINARY)-$(GOOS)-$(GOARCH) $(SRC_LOCATION)

.PHONY: go-run
go-run:
	@go run $(SRC_LOCATION)

.PHONY: go-version
go-version:
	@go run $(SRC_LOCATION) --version

.PHONY: clean
clean:
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
