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
	@which npm 1>/dev/null && cd ui && npm install 1>/dev/null && npm run build 1>/dev/null

install: build ui
	@go install -mod=readonly ./...

cli: build
	@go install -mod=readonly ./...
	
.PHONY: all mod build ui install