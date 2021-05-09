# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=discord-house-cup
COMMIT := $(shell git rev-parse HEAD)
VERSION := "local-dev"

all: lint test
docker: build-docker run-docker
run: build run-local
build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -s -w" -v
lint:
	golangci-lint run
reportcard:
	goreportcard-cli -t 100 -v
test:
	GO111MODULE=on $(GOCMD) test -v --bench --benchmem -coverprofile coverage.txt -covermode=atomic ./...
	GO111MODULE=on $(GOCMD) vet ./... 2> govet-report.out
	GO111MODULE=on $(GOCMD) tool cover -html=coverage.txt -o cover-report.html
	printf "\nCoverage report available at cover-report.html\n\n"
tidy:
	$(GOCMD) mod tidy
clean:
	$(GOCLEAN)
	$(GOCMD) fmt ./...
	rm -f $(BINARY_NAME)
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -s -w" -v
build-docker:
	docker build --build-arg version=$(VERSION) --build-arg commit=$(COMMIT) -t quay.io/sudermanjr/$(BINARY_NAME):dev .
run-docker:
	docker run --rm -p 4004:4004 --env DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN} --env DISCORD_GUILD_ID=${DISCORD_GUILD_ID} quay.io/sudermanjr/$(BINARY_NAME):dev server
run-local:
	go run main.go server --token ${DISCORD_BOT_TOKEN} --guild ${DISCORD_GUILD_ID}
