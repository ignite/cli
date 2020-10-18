---
order: 21
---

# go.mod and Makefile

## `Makefile`

Help users build your application by writing a `./Makefile` in the root directory that includes common commands. The scaffolding tool has created a generic makefile that you will be able to use:

> _*NOTE*_: The below Makefile contains some of same commands as the Cosmos SDK and Tendermint Makefiles.

```
PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=NameService \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=nameserviced \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=nameservicecli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		@echo "--> Installing nameserviced & nameservicecli"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/nameserviced
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/nameservicecli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)
```

### How about including Ledger Nano S support?

This requires a few small changes:

- Create a file `Makefile.ledger` with the following content:

+++ https://github.com/cosmos/sdk-tutorials/blob/master/nameservice/Makefile.ledger

- Add `include Makefile.ledger` at the beginning of the Makefile:

```
BUILD_FLAGS := -ldflags '$(ldflags)'

include Makefile.ledger
all: install
```

## `go.mod`

Golang has a few dependency management tools. In this tutorial you will be using [`Go Modules`](https://github.com/golang/go/wiki/Modules). `Go Modules` uses a `go.mod` file in the root of the repository to define what dependencies the application needs. Cosmos SDK apps currently depend on specific versions of some libraries. The below manifest contains all the necessary versions. To get started replace the contents of the `./go.mod` file with the `constraints` and `overrides` below:

> _*NOTE*_: If you are following along in your own repo you will need to change the module path to reflect that (`github.com/{ .Username }/{ .Project.Repo }`).

- You will have to run `go get ./...` to get all the modules the application is using. This command will get the dependency version stated in the `go.mod` file.
- If you would like to use a specific version of a dependency then you have to run `go get github.com/<github_org>/<repo_name>@<version>`

<!-- <<< @/nameservice/go.mod -->

## Building the app

```bash
# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
nameserviced help
nameservicecli help
```

### Congratulations, you have finished your nameservice application! Try running and interacting with it!
