DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')

VERSION = $(shell git describe --tags)
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = -X github.com/tendermint/starport/starport/internal/version.Version='$(VERSION)' \
	-X github.com/tendermint/starport/starport/internal/version.Head='$(HEAD)' \
	-X github.com/tendermint/starport/starport/internal/version.Date='$(DATE)'
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'

all: install

mod:
	@go mod tidy

pre-build:
	@echo "Fetching latest tags"
	@git fetch --tags

build: mod pre-build
	@go get -u github.com/gobuffalo/packr/v2/packr2
	@cd ./starport/interface/cli/starport && packr2
	@mkdir -p build/
	@go-bindata -pkg cosmosfaucet -prefix starport/pkg/cosmosfaucet -o starport/pkg/cosmosfaucet/openapi_generated.go starport/pkg/cosmosfaucet/openapi/...
	@go build $(BUILD_FLAGS) -o build/ ./starport/interface/cli/...
	@packr2 clean
	@go mod tidy

clean:
	@rm -rf build

ui:
	@rm -rf starport/ui/dist
	-@which npm 1>/dev/null && cd starport/ui && npm install 1>/dev/null && npm run build 1>/dev/null
	go get github.com/rakyll/statik

install: ui build
	@go install $(BUILD_FLAGS) ./...

cli: build
	@go install $(BUILD_FLAGS) ./...

lint:
	golangci-lint run --out-format=tab --issues-exit-code=0
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s

.PHONY: lint

.PHONY: all mod pre-build build ui install
