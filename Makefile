#! /usr/bin/make -f

# Project variables.
PROJECT_NAME = ignite
DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = 
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'
BUILD_FOLDER = ./dist

## install: Install de binary.
install:
	@echo Installing Ignite CLI...
	@go install $(BUILD_FLAGS) ./...
	@ignite version

## build: Build the binary.
build:
	@echo Building Ignite CLI...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@go build $(BUILD_FLAGS) -o $(BUILD_FOLDER) ./...

## prepare snapcraft config for release
snapcraft:
	@sed -i 's/{{version}}/'$(version)'/' packaging/snap/snapcraft.yaml

## mocks: generate mocks
mocks:
	@echo Generating mocks
	@go install github.com/vektra/mockery/v2
	@go generate ./...


## clean: Clean build files. Also runs `go clean` internally.
clean:
	@echo Cleaning build cache...
	@-rm -rf $(BUILD_FOLDER) 2> /dev/null
	@go clean ./...

.PHONY: install build mocks clean

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go vet ./...

## govulncheck: Run govulncheck
govulncheck:
	@echo Running govulncheck...
	@go tool golang.org/x/vuln/cmd/govulncheck ./...

## format: Install and run goimports and gofumpt
format:
	@echo Formatting...
	@go tool mvdan.cc/gofumpt -w .
	@go tool golang.org/x/tools/cmd/goimports -w -local github.com/ignite/cli/v29 .
	@go tool github.com/tbruyelle/mdgofmt/cmd/mdgofmt -w docs

## lint: Run Golang CI Lint.
lint:
	@echo Running golangci-lint...
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run --out-format=tab --issues-exit-code=0

lint-fix:
	@echo Running golangci-lint...
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run --fix --out-format=tab --issues-exit-code=0

.PHONY: govet format lint

## proto-all: Format, lint and generate code from proto files using buf.
proto-all: proto-format proto-lint proto-gen format

## proto-gen: Run buf generate.
proto-gen:
	@echo Generating code from proto...
	@buf generate --template ./proto/buf.gen.yaml --output ./

## proto-format: Run buf format and update files with invalid proto format>
proto-format:
	@echo Formatting proto files...
	@buf format --write

## proto-lint: Run buf lint.
proto-lint:
	@echo Linting proto files...
	@buf lint

.PHONY: proto-all proto-gen proto-format proto-lint

## test-unit: Run the unit tests.
test-unit:
	@echo Running unit tests...
	@go test -race -failfast -v ./ignite/...

## test-integration: Run the integration tests.
test-integration: install
	@echo Running integration tests...
	@go test -race -failfast -v -timeout 60m ./integration/...

## test: Run unit and integration tests.
test: govet govulncheck test-unit test-integration

.PHONY: test-unit test-integration test

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: help

.DEFAULT_GOAL := install
