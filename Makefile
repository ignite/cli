#! /usr/bin/make -f

# Project variables.
PROJECT_NAME = starport
DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = -X github.com/tendermint/starport/starport/internal/version.Head='$(HEAD)' \
	-X github.com/tendermint/starport/starport/internal/version.Date='$(DATE)'
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'
BUILD_FOLDER = ./dist

## install: Install de binary.
install:
	@echo Installing Starport...
	@go install $(BUILD_FLAGS) ./...
	@starport version

## build: Build the binary.
build:
	@echo Building Starport...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@go build $(BUILD_FLAGS) -o $(BUILD_FOLDER) ./...

## clean: Clean build files. Also runs `go clean` internally.
clean:
	@echo Cleaning build cache...
	@-rm -rf $(BUILD_FOLDER) 2> /dev/null
	@go clean ./...

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go vet ./...

## format: Run gofmt.
format:
	@echo Formatting...
	@find . -name '*.go' -type f | xargs gofmt -d -s

## lint: Run Golang CI Lint.
lint:
	@echo Running gocilint...
	@golangci-lint run --out-format=tab --issues-exit-code=0

## test-unit: Run the unit tests.
test-unit:
	@echo Running unit tests...
	@go test -race -failfast -v ./starport/...

## test-integration: Run the integration tests.
test-integration: install
	@echo Running integration tests...
	@go test -race -failfast -v -timeout 60m ./integration/...

## test: Run unit and integration tests.
test: govet test-unit test-integration

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.DEFAULT_GOAL := install
