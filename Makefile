#! /usr/bin/make -f

# Project variables.
VERSION := development
PROJECT_NAME := $(shell basename "$(PWD)")
DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = -X github.com/tendermint/starport/starport/internal/version.Version='$(VERSION)' \
	-X github.com/tendermint/starport/starport/internal/version.Head='$(HEAD)' \
	-X github.com/tendermint/starport/starport/internal/version.Date='$(DATE)'
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOCMD := $(GOBASE)/starport/cmd/starport
GOCILINT := $(GOPATH)/bin/golangci-lint
GORELEASER := $(GOPATH)/bin/goreleaser
STARPORT := $(GOPATH)/bin/starport

# For go 1.13-1.15 compatibility
# https://maelvls.dev/go111module-everywhere/
export GO111MODULE=on

## install: Install de binary.
install:
	@echo Installing Starport...
	@go install $(BUILD_FLAGS) ./...
	@starport version

## build-binary: Build the binary.
build-binary: clean
	@echo Building Starport...
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(PROJECT_NAME) $(GOCMD)

$(STARPORT): install

$(GORELEASER):
	@echo Installing goreleaser...
	@go install github.com/goreleaser/goreleaser@latest

## goreleaser: Build the binaries with Goreleaser.
goreleaser: clean $(GORELEASER)
	@echo Building with goreleaser...
	@$(GORELEASER) --snapshot --skip-publish

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm -r $(GOBIN) 2> /dev/null
	@echo Cleaning build cache...
	@go clean ./...

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go vet ./...

## format: Run gofmt.
format:
	@echo Formatting gocilint...
	@find . -name '*.go' -type f | xargs gofmt -d -s

$(GOCILINT):
	@echo Installing gocilint...
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## lint: Run Golang CI Lint.
lint: $(GOCILINT)
	@echo Running gocilint...
	@$(GOCILINT) run --out-format=tab --issues-exit-code=0

## unit-test: Run the unit tests.
unit-test:
	@echo Running unit tests...
	@go test -race -failfast -v ./starport/...

## integration-test: Run the integration tests.
integration-test: $(STARPORT)
	@echo Running integration tests...
	@go test -race -failfast -v -timeout 60m ./integration/...

## test: Run unit and integration tests.
test: govet unit-test integration-test

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.DEFAULT_GOAL := install
