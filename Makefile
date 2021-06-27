# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=riley
COMMIT := $(shell git rev-parse HEAD)
VERSION := "local-dev"

all: lint test
docker: build-docker run-docker
build:
	$(GOBUILD) cmd/main.go -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -s -w" -v
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
run-dev:
	docker stop dev-postgres 2>/dev/null || true
	docker run --rm -itd -p 5432:5432 --name dev-postgres --env TZ=America/Denver --env POSTGRES_PASSWORD=test --env POSTGRES_DB=test postgres:12
