GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION?=$(shell git describe --abbrev=0 --tags)
BINARY?=bin/bashbot
SRC_LOCATION?=cmd/bashbot/bashbot.go
LDFLAGS="-X main.Version=${VERSION}"
GO_BUILD=go build -ldflags=$(LDFLAGS)
.PHONY: cross
cross:
	rm -rf $(BINARY)*
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-linux-amd64 $(SRC_LOCATION)
	GOOS=linux   GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-linux-arm64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) -o $(BINARY)-darwin-amd64 $(SRC_LOCATION)
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) -o $(BINARY)-darwin-arm64 $(SRC_LOCATION)

.PHONY: build
build:
	rm -rf $(BINARY)-$(GOOS)-$(GOARCH)
	$(GO_BUILD) -o $(BINARY)-$(GOOS)-$(GOARCH) $(SRC_LOCATION)

.PHONY: run
run:
	@go run $(SRC_LOCATION) --config-file $(PWD)/config.json --slack-token $(SLACK_TOKEN)

.PHONY: clean
clean:
	echo "Removing any existing go-binaries"
	rm -rf $(BINARY)*