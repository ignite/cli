all: install

mod:
	@go mod tidy

build: mod
	@go get -u github.com/gobuffalo/packr/v2/packr2
	@packr2
	@go build -mod=readonly -o build/starport main.go
	@packr2 clean
	@go mod tidy

ui:
	@rm -rf ui/dist
	-@which npm 1>/dev/null && cd ui && npm install 1>/dev/null && npm run build 1>/dev/null

install: ui build
	@go install -mod=readonly ./...

cli: build
	@go install -mod=readonly ./...

lint:
	golangci-lint run --out-format=tab --issues-exit-code=0
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
.PHONY: lint
	
.PHONY: all mod build ui install
