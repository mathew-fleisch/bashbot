TARGET := $(shell echo $${PWD\#\#*/})
PKGS := $(shell go list ./... | grep -v /vendor)
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GIT_COMMIT := $(shell git rev-parse HEAD)

LDFLAGS=-ldflags "-X main.GitCommit=$(GIT_COMMIT)"

$(TARGET): $(SRC)
	@echo "compiling $(TARGET)"
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@echo "Cleaning..."
	@rm -f $(TARGET)

install:
	@echo "Installing..."
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

run: install
	@$(TARGET)

docker-no-sudo:
	go mod vendor
	docker build -t testbashbot .

docker:
	go mod vendor
	docker build -t testbashbot .

docker-run: docker
	docker run testbashbot:latest

test:
	go test -cover $(PKGS)

.PHONY: all build clean install uninstall docker docker-run test run

