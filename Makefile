SHELL := /bin/bash

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)


PROTOC_VERSION := 3.12.1
PB_REL := "https://github.com/protocolbuffers/protobuf/releases"

ifeq ($(UNAME_OS),Linux)
  PROTOC_ZIP := protoc-${PROTOC_VERSION}-linux-x86_64.zip
endif
ifeq ($(UNAME_OS),Darwin)
  PROTOC_ZIP := protoc-${PROTOC_VERSION}-osx-x86_64.zip
endif

all: protoc-install install

mod:
	@go mod tidy

build: mod
	@go get -u github.com/gobuffalo/packr/v2/packr2
	@packr2
	@mkdir -p build/
	@go build -mod=readonly -o build/ ./starport/interface/cli/...
	@packr2 clean
	@go mod tidy

ui:
	@rm -rf starport/ui/dist
	-@which npm 1>/dev/null && cd starport/ui && npm install 1>/dev/null && npm run build 1>/dev/null

install: ui build
	@go install -mod=readonly ./...

cli: build
	@go install -mod=readonly ./...

lint:
	golangci-lint run --out-format=tab --issues-exit-code=0
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s

protoc-install:
ifeq (, $(shell which protoc))
	@echo "installing protoc..."
	echo "export PATH=${PATH}:${HOME}/bin" >> ~/.bash_profile
	source ~/.bash_profile
	curl -LOs ${PB_REL}/download/v${PROTOC_VERSION}/${PROTOC_ZIP}
	unzip ${PROTOC_ZIP} bin/protoc 
	mkdir -p ${HOME}/bin
	mv bin/protoc ${HOME}/bin
	rm ${PROTOC_ZIP}
	protoc --version
endif

.PHONY: lint
	
.PHONY: all mod build ui install
