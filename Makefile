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
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -s -w" -v cmd/root/main.go
lint:
	golangci-lint run
reportcard:
	goreportcard-cli -t 100 -v
test: dev-db-start
	$(GOCMD) test -v --bench --benchmem -coverprofile coverage.txt -covermode=atomic ./... -args -env dev
	$(GOCMD) vet ./... 2> govet-report.out
	$(GOCMD) tool cover -html=coverage.txt -o cover-report.html
	printf "\nCoverage report available at cover-report.html\n\n"
tidy:
	$(GOCMD) mod tidy
clean:
	$(GOCLEAN)
	$(GOCMD) fmt ./...
	rm -f $(BINARY_NAME)
dev-db-start:
	docker run --rm -itd -p 5432:5432 --name riley-test-postgres --env TZ=America/Denver --env POSTGRES_PASSWORD=test --env POSTGRES_DB=test postgres:12 || true
	go run cmd/root/main.go db init
dev-db-stop:
	docker stop riley-test-postgres
