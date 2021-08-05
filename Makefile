GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION?=$(shell git describe --abbrev=0 --tags)
LATEST_VERSION?=$(shell curl -s https://api.github.com/repos/mathew-fleisch/bashbot/releases/latest | grep tag_name | cut -d '"' -f 4)
BINARY?=bin/bashbot
SRC_LOCATION?=cmd/bashbot/bashbot.go
LDFLAGS="-X main.Version=${VERSION}"
GO_BUILD=go build -ldflags=$(LDFLAGS)

.PHONY: setup
setup:
	go mod tidy
	go mod vendor
	go install -v ./...
	go get github.com/slack-go/slack@master
	go get github.com/sirupsen/logrus

.PHONY: cross
cross:
	rm -rf $(BINARY)*
	go mod tidy
	go mod vendor
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-linux-amd64 $(SRC_LOCATION)
	GOOS=linux   GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-linux-arm64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-darwin-amd64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-darwin-arm64 $(SRC_LOCATION)

.PHONY: build
build:
	rm -rf $(BINARY)-$(GOOS)-$(GOARCH)
	go mod tidy
	go mod vendor
	$(GO_BUILD) -o $(BINARY)-$(GOOS)-$(GOARCH) $(SRC_LOCATION)

.PHONY: run
run:
	@go run $(SRC_LOCATION) --config-file $(PWD)/config.json --slack-token $(SLACK_TOKEN)

.PHONY: run-version
run:
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
	bashbot --help