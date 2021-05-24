DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
VERSION = $(shell git describe --tags)
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = -X github.com/tendermint/starport/starport/internal/version.Version='$(VERSION)' \
	-X github.com/tendermint/starport/starport/internal/version.Head='$(HEAD)' \
	-X github.com/tendermint/starport/starport/internal/version.Date='$(DATE)'
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'

pre-build:
	@git fetch --tags
	@./scripts/gen-nodetime

install: pre-build
	@echo Installing Starport...
	@go install $(BUILD_FLAGS) ./...
	@starport version

format:
	@find . -name '*.go' -type f | xargs gofmt -d -s

lint:
	@golangci-lint run --out-format=tab --issues-exit-code=0

ui:
	@rm -rf starport/ui/app/dist
	-@which npm 1>/dev/null && cd starport/ui/app && npm install 1>/dev/null && npm run build 1>/dev/null

.DEFAULT_GOAL := install 
